package encrypt

func EncryptMap(data map[string]any, key string) (map[string]any, error) {
	encrypted := make(map[string]any)
	for k, v := range data {
		if x, ok := v.(string); ok {
			s, err := EncryptString(x, key)
			if err != nil {
				return nil, err
			}
			encrypted[k] = s
		}
	}
	return encrypted, nil
}

func DecryptMap(encryptedData map[string]any, key string) (map[string]any, error) {
	decrypted := make(map[string]any)
	for k, v := range encryptedData {
		if x, ok := v.(string); ok {
			s, err := DecryptString(x, key)
			if err != nil {
				return nil, err
			}
			decrypted[k] = s
		}
	}
	return decrypted, nil
}
