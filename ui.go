package main

import (
	"fmt"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"

	"gioui.org/io/clipboard"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var padding1 = unit.Dp(8)
var padding2 = unit.Dp(4)
var padding3 = unit.Dp(2)

type Application struct {
	Theme *material.Theme

	IPv4DecHexBinConverter    IPv4DecHexBinConverter
	NetMaskCIDRSlashConverter NetMaskCIDRSlashConverter
	NetAddrFinder             NetAddrFinder
	IPInfoChecker             IPInfoChecker
	DecHexBinConverter        DecHexBinConverter
	ANDOperationOnTwoBins     ANDOperationOnTwoBins
}

func NewApplication() *Application {
	theme := material.NewTheme(gofont.Collection())
	theme.TextSize = 12
	application := Application{
		Theme: theme,
	}

	return &application
}

func (a *Application) Run(window *app.Window) error {
	var ops op.Ops

	for e := range window.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			a.Layout(gtx)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}

func (a *Application) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(padding2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(padding3).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			spacer := layout.Rigid(layout.Spacer{Height: padding1}.Layout)

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(Heading(a.Theme, "Convert IPv4 dotted decimal to hexadecimal and binary:").Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.IPv4DecHexBinConverter.Layout(a.Theme, gtx)
				}),
				spacer,
				layout.Rigid(Heading(a.Theme, "Convert between network mask and CIDR slash value:").Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.NetMaskCIDRSlashConverter.Layout(a.Theme, gtx)
				}),
				spacer,
				layout.Rigid(Heading(a.Theme, "Compute network address from host IP address and network mask:").Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.NetAddrFinder.Layout(a.Theme, gtx)
				}),
				spacer,
				layout.Rigid(Heading(a.Theme, "Is IP address private/loopback/link-local unicast/multicast?").Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.IPInfoChecker.Layout(a.Theme, gtx)
				}),
				spacer,
				layout.Rigid(Heading(a.Theme, "Convert between decimal, hexadecimal, and binary formats:").Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.DecHexBinConverter.Layout(a.Theme, gtx)
				}),
				spacer,
				layout.Rigid(Heading(a.Theme, "Perform AND operation on two binary numbers:").Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.ANDOperationOnTwoBins.Layout(a.Theme, gtx)
				}),
			)
		})
	})
}

type IPv4DecHexBinConverter struct {
	Dec Field
	Hex Field
	Bin Field
}

func (conv *IPv4DecHexBinConverter) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if conv.Dec.Changed() {
		hexValue, err := IPv4ToHexFormat(conv.Dec.Text())
		conv.Dec.Invalid = err != nil
		conv.Hex.SetText(hexValue)

		binValue, err := IPv4ToBinFormat(conv.Dec.Text())
		conv.Dec.Invalid = err != nil
		conv.Bin.SetText(FormatBinInNimbles(binValue))
	}

	if conv.Hex.Changed() {
		ipv4Value, err := HexToIPv4Format(conv.Hex.Text())
		conv.Hex.Invalid = err != nil

		if conv.Hex.Invalid {
			conv.Dec.SetText("")
			conv.Bin.SetText("")
		} else {
			conv.Dec.SetText(ipv4Value)

			binValue, _ := IPv4ToBinFormat(ipv4Value)
			conv.Bin.SetText(binValue)
		}
	}

	if conv.Bin.Changed() {
		ipv4Value, err := BinToIPv4Format(conv.Bin.Text())
		conv.Bin.Invalid = err != nil

		if conv.Bin.Invalid {
			conv.Dec.SetText("")
			conv.Hex.SetText("")
		} else {
			conv.Dec.SetText(ipv4Value)

			hexValue, _ := IPv4ToHexFormat(ipv4Value)
			conv.Hex.SetText(hexValue)
		}
	}

	spacer := layout.Rigid(layout.Spacer{Width: padding2}.Layout)

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(th, "Dec:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Dec.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Hex:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Hex.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Bin:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Bin.Layout(th, gtx)
		}),
	)
}

type NetMaskCIDRSlashConverter struct {
	NetMask   Field
	CIDRSlash Field
}

func (conv *NetMaskCIDRSlashConverter) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if conv.NetMask.Changed() {
		cidrSlashValue, err := NetworkMaskToCIDRSlashValue(conv.NetMask.Text())
		conv.NetMask.Invalid = err != nil
		conv.CIDRSlash.SetText(cidrSlashValue)
	}

	if conv.CIDRSlash.Changed() {
		netMaskValue, err := CIDRSlashValueToNetworkMask(conv.CIDRSlash.Text())
		conv.CIDRSlash.Invalid = err != nil
		conv.NetMask.SetText(netMaskValue)
	}

	spacer := layout.Rigid(layout.Spacer{Width: padding2}.Layout)

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(th, "Network mask:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.NetMask.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "CIDR slash value:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.CIDRSlash.Layout(th, gtx)
		}),
	)
}

type NetAddrFinder struct {
	HostIP       Field
	NetMask      Field
	NetAddr      widget.Clickable
	NetAddrValue string
}

func (finder *NetAddrFinder) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if finder.HostIP.Changed() || finder.NetMask.Changed() {
		var err error

		if finder.HostIP.Text() != "" && finder.NetMask.Text() != "" {
			finder.NetAddrValue, err = FindNetworkAddress(finder.HostIP.Text(), finder.NetMask.Text())
			finder.HostIP.Invalid = err != nil
			finder.NetMask.Invalid = err != nil
		}
	}

	if finder.NetAddr.Clicked() {
		clipboard.WriteOp{Text: finder.NetAddrValue}.Add(gtx.Ops)
	}

	spacer := layout.Rigid(layout.Spacer{Width: padding2}.Layout)

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(th, "Host IP:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return finder.HostIP.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Network mask:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return finder.NetMask.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Network address:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &finder.NetAddr, material.Body1(th, finder.NetAddrValue).Layout)
		}),
	)
}

type IPInfoChecker struct {
	IPAddr                       Field
	PrivateChecked               widget.Clickable
	PrivateCheckedValue          bool
	LoopbackChecked              widget.Clickable
	LoopbackCheckedValue         bool
	LinkLocalUnicastChecked      widget.Clickable
	LinkLocalUnicastCheckedValue bool
	MulticastChecked             widget.Clickable
	MulticastCheckedValue        bool
}

func (checker *IPInfoChecker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	var privateChecked, loopbackChecked, linklocalUnicastChecked, multicastChecked string

	if checker.IPAddr.Changed() {
		var err error

		checker.PrivateCheckedValue, err = IsPrivateIP(checker.IPAddr.Text())
		checker.IPAddr.Invalid = err != nil

		checker.LoopbackCheckedValue, err = IsLoopbackIP(checker.IPAddr.Text())
		checker.IPAddr.Invalid = err != nil

		checker.LinkLocalUnicastCheckedValue, err = IsLinkLocalUnicastIP(checker.IPAddr.Text())
		checker.IPAddr.Invalid = err != nil

		checker.MulticastCheckedValue, err = IsMulticastIP(checker.IPAddr.Text())
		checker.IPAddr.Invalid = err != nil
	}

	if checker.IPAddr.Text() != "" {
		privateChecked = strconv.FormatBool(checker.PrivateCheckedValue)
		loopbackChecked = strconv.FormatBool(checker.LoopbackCheckedValue)
		linklocalUnicastChecked = strconv.FormatBool(checker.LinkLocalUnicastCheckedValue)
		multicastChecked = strconv.FormatBool(checker.MulticastCheckedValue)
	}

	if checker.PrivateChecked.Clicked() {
		clipboard.WriteOp{Text: privateChecked}.Add(gtx.Ops)
	} else if checker.LoopbackChecked.Clicked() {
		clipboard.WriteOp{Text: loopbackChecked}.Add(gtx.Ops)
	} else if checker.LinkLocalUnicastChecked.Clicked() {
		clipboard.WriteOp{Text: linklocalUnicastChecked}.Add(gtx.Ops)
	} else if checker.MulticastChecked.Clicked() {
		clipboard.WriteOp{Text: multicastChecked}.Add(gtx.Ops)
	}

	spacer := layout.Rigid(layout.Spacer{Width: padding2}.Layout)

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(th, "IP:").Layout),
		spacer,
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			return checker.IPAddr.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Private?").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &checker.PrivateChecked, material.Body1(th, privateChecked).Layout)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Loopback?").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &checker.LoopbackChecked, material.Body1(th, loopbackChecked).Layout)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Link-local unicast?").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &checker.LinkLocalUnicastChecked, material.Body1(th, linklocalUnicastChecked).Layout)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Multicast?").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &checker.MulticastChecked, material.Body1(th, multicastChecked).Layout)
		}),
	)
}

type DecHexBinConverter struct {
	Dec Field
	Hex Field
	Bin Field
}

func (conv *DecHexBinConverter) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if conv.Dec.Changed() {
		trimmedDec := strings.TrimSpace(conv.Dec.Text())

		hexValue, err := DecToHex(trimmedDec)
		conv.Dec.Invalid = err != nil
		conv.Hex.SetText(hexValue)

		binValue, err := DecToBin(trimmedDec)
		conv.Dec.Invalid = err != nil
		conv.Bin.SetText(FormatBinInNimbles(binValue))
	}

	if conv.Hex.Changed() {
		decValue, err := strconv.ParseInt(strings.ReplaceAll(conv.Hex.Text(), " ", ""), 16, 64)
		conv.Hex.Invalid = err != nil

		if conv.Hex.Invalid {
			conv.Dec.SetText("")
			conv.Bin.SetText("")
		} else {
			conv.Dec.SetText(fmt.Sprintf("%d", decValue))
			conv.Bin.SetText(FormatBinInNimbles(fmt.Sprintf("%b", decValue)))
		}
	}

	if conv.Bin.Changed() {
		decValue, err := strconv.ParseInt(strings.ReplaceAll(conv.Bin.Text(), " ", ""), 2, 64)
		conv.Bin.Invalid = err != nil

		if conv.Bin.Invalid {
			conv.Dec.SetText("")
			conv.Hex.SetText("")
		} else {
			conv.Dec.SetText(fmt.Sprintf("%d", decValue))
			conv.Hex.SetText(fmt.Sprintf("%X", decValue))
		}
	}

	spacer := layout.Rigid(layout.Spacer{Width: padding2}.Layout)

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(th, "Dec:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Dec.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Hex:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Hex.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Bin:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Bin.Layout(th, gtx)
		}),
	)
}

type ANDOperationOnTwoBins struct {
	Bin1        Field
	Bin2        Field
	Result      widget.Clickable
	ResultValue string
}

func (conv *ANDOperationOnTwoBins) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	var err error

	var dec1, dec2 int64

	if conv.Bin1.Changed() || conv.Bin2.Changed() {
		conv.ResultValue = ""
		trimmedBin1 := strings.ReplaceAll(conv.Bin1.Text(), " ", "")
		trimmedBin2 := strings.ReplaceAll(conv.Bin2.Text(), " ", "")

		dec1, err = strconv.ParseInt(trimmedBin1, 2, 64)
		conv.Bin1.Invalid = err != nil

		dec2, err = strconv.ParseInt(trimmedBin2, 2, 64)
		conv.Bin2.Invalid = err != nil

		if conv.Bin1.Text() != "" && conv.Bin2.Text() != "" {
			if !conv.Bin1.Invalid && !conv.Bin2.Invalid {
				maxBinLength := len(trimmedBin1)
				if len(trimmedBin2) > len(trimmedBin1) {
					maxBinLength = len(trimmedBin2)
				}

				binNum := fmt.Sprintf(fmt.Sprintf("%%0%db", maxBinLength), dec1&dec2)
				conv.ResultValue = FormatBinInNimbles(binNum)
			}
		}
	}

	if conv.Result.Clicked() {
		clipboard.WriteOp{Text: conv.ResultValue}.Add(gtx.Ops)
	}

	spacer := layout.Rigid(layout.Spacer{Width: padding2}.Layout)

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(th, "First:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Bin1.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Second:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Bin2.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Result:").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &conv.Result, material.Body1(th, conv.ResultValue).Layout)
		}),
	)
}
