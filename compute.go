package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func IPv4ToHexFormat(ipAddress string) (string, error) {
	ipv4 := net.ParseIP(ipAddress)
	if ipv4 == nil {
		return "", errors.New("IPv4ToHexFormat: ipAddress format is invalid")
	}

	if ipv4.To4() == nil {
		return "", errors.New("IPv4ToHexFormat: ipAddress is not IPv4")
	}

	return fmt.Sprintf("%02X%02X%02X%02X", ipv4[12], ipv4[13], ipv4[14], ipv4[15]), nil
}

func IPv4ToBinFormat(ipAddress string) (string, error) {
	ipv4 := net.ParseIP(ipAddress)
	if ipv4 == nil {
		return "", errors.New("IPv4ToBinFormat: ipAddress format is invalid")
	}

	if ipv4.To4() == nil {
		return "", errors.New("IPv4ToBinFormat: ipAddress is not IPv4")
	}

	var buf bytes.Buffer

	for i, octet := range ipv4[12:] {
		octetStr := fmt.Sprintf("%08b", octet)
		buf.WriteString(octetStr[:4])
		buf.WriteString(" ")
		buf.WriteString(octetStr[4:])

		if i < 3 {
			buf.WriteString(" ")
		}
	}

	return buf.String(), nil
}

func BinToIPv4Format(binNumber string) (string, error) {
	trimmedBinNumber := strings.ReplaceAll(binNumber, " ", "")
	if len(trimmedBinNumber) != 32 {
		return "", errors.New("BinToIPv4Format: binNumber is invalid")
	}

	var ipv4 bytes.Buffer

	for i := 0; i < 32; i += 8 {
		decValue, err := strconv.ParseInt(trimmedBinNumber[i:i+8], 2, 64)
		if err != nil {
			return "", errors.New("BinToIPv4Format: binNumber is invalid")
		}

		ipv4.WriteString(fmt.Sprintf("%d", decValue))

		if i < 24 {
			ipv4.WriteString(".")
		}
	}

	return ipv4.String(), nil
}

func HexToIPv4Format(hexNumber string) (string, error) {
	trimmedHexNumber := strings.ReplaceAll(hexNumber, " ", "")
	if len(trimmedHexNumber) != 8 {
		return "", errors.New("HexToIPv4Format: hexNumber is invalid")
	}

	var ipv4 bytes.Buffer

	for i := 0; i < 8; i += 2 {
		decValue, err := strconv.ParseInt(trimmedHexNumber[i:i+2], 16, 64)
		if err != nil {
			return "", errors.New("HexToIPv4Format: hexNumber is invalid")
		}

		ipv4.WriteString(fmt.Sprintf("%d", decValue))

		if i < 6 {
			ipv4.WriteString(".")
		}
	}

	return ipv4.String(), nil
}

func NetworkMaskToCIDRSlashValue(netMask string) (string, error) {
	ipv4 := net.ParseIP(netMask)
	if ipv4 == nil {
		return "", errors.New("NetworkMaskToCIDRSlashValue: ipAddress format is invalid")
	}

	if ipv4.To4() == nil {
		return "", errors.New("NetworkMaskToCIDRSlashValue: ipAddress is not IPv4")
	}

	mask := net.IPv4Mask(ipv4[12], ipv4[13], ipv4[14], ipv4[15])
	ones, _ := mask.Size()

	return "/" + strconv.Itoa(ones), nil
}

func CIDRSlashValueToNetworkMask(cidrSlashValue string) (string, error) {
	if cidrSlashValue == "" {
		return "", errors.New("CIDRSlashValueToNetworkMask: cidrSlashValue is empty")
	}

	if cidrSlashValue[0] != '/' {
		return "", errors.New("CIDRSlashValueToNetworkMask: cidrSlashValue format is invalid")
	}

	ones, err := strconv.Atoi(cidrSlashValue[1:])
	if err != nil {
		return "", errors.New("CIDRSlashValueToNetworkMask: cidrSlashValue is missing the number part")
	}

	if ones > 32 {
		return "", errors.New("CIDRSlashValueToNetworkMask: cidrSlashValue cannot be bigger than 32")
	}

	if mask := net.CIDRMask(ones, 32); mask != nil {
		return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]), nil
	} else {
		return "", errors.New("CIDRSlashValueToNetworkMask: cidrSlashValue format is invalid")
	}
}

func FindNetworkAddress(hostIPAddress string, networkMask string) (string, error) {
	cidrSlashValue, err := NetworkMaskToCIDRSlashValue(networkMask)
	if err != nil {
		return "", fmt.Errorf("FindNetworkAddress: %w", err)
	}

	_, ipNet, err := net.ParseCIDR(hostIPAddress + cidrSlashValue)
	if err != nil {
		return "", fmt.Errorf("FindNetworkAddress: %w", err)
	}

	return ipNet.String(), err
}

func IsPrivateIP(ipAddress string) (bool, error) {
	if ipAddr := net.ParseIP(ipAddress); ipAddr == nil {
		return false, errors.New("IsPrivateIP: ipAddress is invalid")
	} else {
		return ipAddr.IsPrivate(), nil
	}
}

func IsLoopbackIP(ipAddress string) (bool, error) {
	if ipAddr := net.ParseIP(ipAddress); ipAddr == nil {
		return false, errors.New("IsLoopbackIP: ipAddress is invalid")
	} else {
		return ipAddr.IsLoopback(), nil
	}
}

func IsLinkLocalUnicastIP(ipAddress string) (bool, error) {
	if ipAddr := net.ParseIP(ipAddress); ipAddr == nil {
		return false, errors.New("IsLinkLocalUnicastIP: ipAddress is invalid")
	} else {
		return ipAddr.IsLinkLocalUnicast(), nil
	}
}

func IsMulticastIP(ipAddress string) (bool, error) {
	if ipAddr := net.ParseIP(ipAddress); ipAddr == nil {
		return false, errors.New("IsMulticastIP: ipAddress is invalid")
	} else {
		return ipAddr.IsMulticast(), nil
	}
}

func DecToBin(decimalNumber string) (string, error) {
	if binNumber, err := strconv.Atoi(decimalNumber); err != nil {
		return "", errors.New("DecToBin: decimalNumber is invalid")
	} else {
		return fmt.Sprintf("%b", binNumber), nil
	}
}

func DecToHex(decimalNumber string) (string, error) {
	if hexNumber, err := strconv.Atoi(decimalNumber); err != nil {
		return "", errors.New("DecToHex: decimalNumber is invalid")
	} else {
		return fmt.Sprintf("%X", hexNumber), nil
	}
}

func FormatBinInNimbles(binNumber string) string {
	var buf bytes.Buffer

	trimmedBinNumber := strings.ReplaceAll(binNumber, " ", "")
	leftMostNimbleEndIndex := len(trimmedBinNumber) % 4

	if leftMostNimbleEndIndex != 0 {
		buf.WriteString(trimmedBinNumber[:leftMostNimbleEndIndex])

		if leftMostNimbleEndIndex+1 < len(trimmedBinNumber) {
			buf.WriteByte(' ')
		}
	}

	for i, c := range trimmedBinNumber[leftMostNimbleEndIndex:] {
		if i != 0 && i%4 == 0 {
			buf.WriteByte(' ')
		}

		buf.WriteRune(c)
	}

	return buf.String()
}
