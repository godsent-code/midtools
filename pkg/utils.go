package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

func WriteResponse(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Error encoding response")
	}
}

func ValidateGhanaLicensePlate(plate string) (bool, string) {
	// Clean the input: remove extra spaces and convert to uppercase
	plate = strings.ToUpper(strings.TrimSpace(plate))
	if plate == "" {
		return false, "Empty plate number"
	}

	// Remove all spaces and hyphens for consistent validation
	cleanPlate := strings.ReplaceAll(strings.ReplaceAll(plate, " ", ""), "-", "")

	// Define valid region codes
	validRegionCodes := map[string]bool{
		// Ashanti Region
		"AC": true, "AE": true, "AK": true, "AP": true, "AS": true, "AW": true,
		// Bono Region
		"BA": true, "BR": true, "BW": true,
		// Bono East
		"BE": true, "BT": true,
		// Central Region
		"CR": true, "CW": true, "CS": true,
		// Eastern Region
		"EN": true, "ER": true, "ES": true,
		// Greater Accra
		"GB": true, "GC": true, "GE": true, "GG": true, "GH": true, "GL": true,
		"GM": true, "GN": true, "GR": true, "GT": true, "GS": true, "GW": true,
		"GX": true, "GY": true,
		// Northern Region
		"NR": true, "NW": true,
		// Upper East
		"UE": true, "UW": true, "UD": true,
		// Upper West
		"UH": true,
		// Volta Region
		"VA": true, "VD": true, "VR": true,
		// Western Region
		"WR": true, "WT": true,
	}

	// Special plates (government, trade, etc.)
	specialCodes := map[string]bool{
		"GA":  true, // Armed Forces
		"GP":  true, // Police
		"FS":  true, // Fire Service
		"PS":  true, // Prisons Service
		"FZB": true, // Free Zone Board
	}

	// ============ PATTERN 1: DV PLATES (WITH YEAR SUFFIX) ============
	// Format: DV + 4 digits + year suffix (e.g., DV123422 or DV1234-22)
	dvPattern := regexp.MustCompile(`^(DV)(\d{4})(\d{2})$`)
	if matches := dvPattern.FindStringSubmatch(cleanPlate); matches != nil {
		return true, fmt.Sprintf("Valid DV trade plate (%s, year 20%s)",
			matches[1]+matches[2], matches[3])
	}

	// ============ PATTERN 2: MOTORCYCLE PLATES ============
	// Format: M + 5 digits (e.g., M12345) or M-12345
	// Note: Sometimes it's M + 4 digits, but 5 is common [citation:1]
	motorcyclePattern1 := regexp.MustCompile(`^M(\d{4,5})$`)
	motorcyclePattern2 := regexp.MustCompile(`^M-(\d{4,5})$`) // With hyphen

	if matches := motorcyclePattern1.FindStringSubmatch(cleanPlate); matches != nil {
		return true, fmt.Sprintf("Valid motorcycle plate (blue background, number %s)", matches[1])
	}
	if matches := motorcyclePattern2.FindStringSubmatch(cleanPlate); matches != nil {
		return true, fmt.Sprintf("Valid motorcycle plate (blue background, number %s)", matches[1])
	}

	// Also check with hyphen (using original plate for hyphen variation)
	if strings.HasPrefix(plate, "M-") || strings.HasPrefix(plate, "M ") {
		numberPart := strings.TrimPrefix(strings.TrimPrefix(plate, "M-"), "M ")
		numberPart = strings.TrimSpace(numberPart)
		if len(numberPart) >= 4 && len(numberPart) <= 5 && regexp.MustCompile(`^\d+$`).MatchString(numberPart) {
			return true, fmt.Sprintf("Valid motorcycle plate (blue background, number %s)", numberPart)
		}
	}

	// ============ PATTERN 3: OLD FORMAT WITH YEAR SUFFIX ============
	// Format: Region + 1-4 digits + year (e.g., GR123422)
	oldFormatPattern := regexp.MustCompile(`^([A-Z]{2})(\d{1,4})(\d{2})$`)
	if matches := oldFormatPattern.FindStringSubmatch(cleanPlate); matches != nil {
		regionCode := matches[1]
		digits := matches[2]
		yearCode := matches[3]

		if validRegionCodes[regionCode] {
			return true, fmt.Sprintf("Valid old format (%s region, digits %s, year 20%s)",
				regionCode, digits, yearCode)
		}
		return false, fmt.Sprintf("Invalid region code: %s", regionCode)
	}

	// ============ PATTERN 4: NEW FORMAT WITH ZONE CODE ============
	// Format: Region + 1-4 digits + zone (e.g., GR1234AD)
	newFormatPattern := regexp.MustCompile(`^([A-Z]{2})(\d{1,4})([A-Z]{2})$`)
	if matches := newFormatPattern.FindStringSubmatch(cleanPlate); matches != nil {
		regionCode := matches[1]
		digits := matches[2]
		zoneCode := matches[3]

		if validRegionCodes[regionCode] {
			return true, fmt.Sprintf("Valid new format (%s region, digits %s, zone %s)",
				regionCode, digits, zoneCode)
		}
		return false, fmt.Sprintf("Invalid region code: %s", regionCode)
	}

	// ============ PATTERN 5: SPECIAL FORMATS ============
	// Format: Special code + digits (e.g., GP5, GA123)
	specialFormatPattern := regexp.MustCompile(`^([A-Z]{2,3})(\d{1,4})$`)
	if matches := specialFormatPattern.FindStringSubmatch(cleanPlate); matches != nil {
		code := matches[1]
		digits := matches[2]

		if specialCodes[code] {
			return true, fmt.Sprintf("Valid special plate (%s, digits %s)", code, digits)
		}

		// Also check if it's a valid region code without year/zone
		if validRegionCodes[code] {
			return true, fmt.Sprintf("Valid basic format (%s region, digits %s)", code, digits)
		}

		return false, fmt.Sprintf("Invalid code: %s", code)
	}

	// ============ PATTERN 6: PERSONALISED PLATES ============
	return validatePersonalisedPlate(plate, cleanPlate)
}

func validatePersonalisedPlate(original string, clean string) (bool, string) {
	// Basic sanity checks for personalised plates
	if len(clean) < 2 || len(clean) > 15 {
		return false, "Does not match any Ghanaian license plate format"
	}

	// Personalised plates can have various patterns
	// Examples: "SERIOUS1-11", "RAPDR1Z", "KEN50-10"

	// Check if it contains letters and numbers
	hasLetters := regexp.MustCompile(`[A-Z]`).MatchString(clean)
	hasNumbers := regexp.MustCompile(`[0-9]`).MatchString(clean)

	if !hasLetters || !hasNumbers {
		return false, "Personalised plates typically contain both letters and numbers"
	}

	// If it passed basic checks, it COULD be a personalised plate
	// In production, you'd need database verification
	return true, "Possible personalised plate (requires DVLA database verification)"
}
