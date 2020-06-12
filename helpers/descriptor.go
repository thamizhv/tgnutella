package helpers

import (
	"fmt"

	"github.com/thamizhv/tgnutella/constants"
	"github.com/thamizhv/tgnutella/models"
)

func GetPayloadDescriptor(descriptorType string) (models.DescriptorType, error) {

	switch descriptorType {
	case constants.PingDescriptor:
		return constants.DescriptorTypePing, nil
	case constants.PongDescriptor:
		return constants.DescriptorTypePong, nil
	case constants.QueryDescriptor:
		return constants.DescriptorTypeQuery, nil
	case constants.QueryHitDescriptor:
		return constants.DescriptorTypeQueryHit, nil
	default:
		return 0, fmt.Errorf("invalid descriptor type: %s", descriptorType)
	}
}

func GetDescriptorString(payloadDescriptor models.DescriptorType) (string, error) {

	switch payloadDescriptor {
	case constants.DescriptorTypePing:
		return constants.PingDescriptor, nil
	case constants.DescriptorTypePong:
		return constants.PongDescriptor, nil
	case constants.DescriptorTypeQuery:
		return constants.QueryDescriptor, nil
	case constants.DescriptorTypeQueryHit:
		return constants.QueryHitDescriptor, nil
	default:
		return "", fmt.Errorf("invalid payload descriptor type: %v", payloadDescriptor)
	}
}
