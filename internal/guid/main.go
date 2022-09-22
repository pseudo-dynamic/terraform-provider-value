package guid

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/uuid"
)

func ComposeGuidSeed(
	providerSeedAddition *string,
	providerMetaSeedAddition *string,
	resourceName string,
	attributeName string,
	resourceGuidSeed *string) string {
	empty := ""

	if providerSeedAddition == nil {
		providerSeedAddition = &empty
	}

	if providerMetaSeedAddition == nil {
		providerMetaSeedAddition = &empty
	}

	if resourceGuidSeed == nil {
		resourceGuidSeed = &empty
	}

	return *providerSeedAddition + "|" +
		*providerMetaSeedAddition + "|" +
		resourceName + "|" +
		attributeName + "|" +
		*resourceGuidSeed
}

func CreateResourceTempDir(resourceName string) (string, error) {
	providerTempDir := filepath.Join(os.TempDir(), "tf-"+resourceName)
	var err error
	if _, err = os.Stat(providerTempDir); os.IsNotExist(err) {
		os.MkdirAll(providerTempDir, 0700) // Create your file
	}
	return providerTempDir, err
}

func GetPlanCachedBoolean(
	isPlanPhase bool,
	composedGuidSeed,
	resourceName string,
	getActualValue func() types.Bool) (types.Bool, error) {
	providerTempDir, _ := CreateResourceTempDir(resourceName)

	deterministicFileName := uuid.DeterministicUuidFromString(composedGuidSeed).String()
	deterministicTempFilePath := filepath.Join(providerTempDir, deterministicFileName)

	var returningValue types.Bool

	if isPlanPhase {
		// This is the plan phase
		var actualValueByte byte
		actualValue := getActualValue()
		if actualValue.IsUnknown() {
			actualValueByte = 2
		} else if actualValue.Value {
			actualValueByte = 1
		}

		deterministicFile, err := os.OpenFile(deterministicTempFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)

		if err == nil {
			defer deterministicFile.Close()
			_, err = deterministicFile.Write([]byte{actualValueByte})
		}

		if err != nil {
			return returningValue, errors.New("error while working with file in the plan phase")
		}

		returningValue = actualValue
	} else if _, err := os.Stat(deterministicTempFilePath); err == nil {
		deterministicFile, err := os.Open(deterministicTempFilePath)
		var readBytes []byte

		if errors.Is(err, os.ErrNotExist) {
			return returningValue, errors.New(`the file does not exist. This can mean
			1. the file got deleted before apply-phase because guid collision or
			2. this plan method got called the third time`)
		} else if err == nil {
			readBytes = make([]byte, 1)
			deterministicFile.Read(readBytes)
			deterministicFile.Close()
		}

		if err != nil {
			return returningValue, errors.New("error while working with the file in the apply phase")
		}

		readByte := readBytes[0]

		if readByte == 2 {
			returningValue = types.Bool{
				Unknown: true,
			}
		} else if readByte == 1 {
			returningValue = types.Bool{
				Value: true,
			}
		} else {
			returningValue = types.Bool{
				Value: false,
			}
		}

		// ISSUE: check why it does not work
		// os.Remove(deterministicTempFilePath)
	} else {
		returningValue = getActualValue()
	}

	return returningValue, nil
}
