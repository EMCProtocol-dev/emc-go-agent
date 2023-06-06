package helper

import (
	"fmt"
	"github.com/emc-go-agent/edge-matrix/crypto"
	"github.com/emc-go-agent/edge-matrix/secrets"
	"github.com/emc-go-agent/edge-matrix/secrets/local"
	"github.com/emc-go-agent/edge-matrix/types"
	"github.com/hashicorp/go-hclog"
)

// SetupLocalSecretsManager is a helper method for boilerplate local secrets manager setup
func SetupLocalSecretsManager(dataDir string) (secrets.SecretsManager, error) {
	return local.SecretsManagerFactory(
		nil, // Local secrets manager doesn't require a config
		&secrets.SecretsManagerParams{
			Logger: hclog.NewNullLogger(),
			Extra: map[string]interface{}{
				secrets.Path: dataDir,
			},
		},
	)
}

// InitECDSAValidatorKey creates new ECDSA key and set as a validator key
func InitECDSAValidatorKey(secretsManager secrets.SecretsManager) (types.Address, error) {
	if secretsManager.HasSecret(secrets.ValidatorKey) {
		return types.ZeroAddress, fmt.Errorf(`secrets "%s" has been already initialized`, secrets.ValidatorKey)
	}

	validatorKey, validatorKeyEncoded, err := crypto.GenerateAndEncodeECDSAPrivateKey()
	if err != nil {
		return types.ZeroAddress, err
	}

	address := crypto.PubKeyToAddress(&validatorKey.PublicKey)

	// Write the validator private key to the secrets manager storage
	if setErr := secretsManager.SetSecret(
		secrets.ValidatorKey,
		validatorKeyEncoded,
	); setErr != nil {
		return types.ZeroAddress, setErr
	}

	return address, nil
}

func InitICPIdentityKey(secretsManager secrets.SecretsManager) ([]byte, error) {
	if secretsManager.HasSecret(secrets.ICPIdentityKey) {
		return nil, fmt.Errorf(`secrets "%s" has been already initialized`, secrets.ICPIdentityKey)
	}
	// generate ed25519 key for ICP identity
	ed25519PubKey, ed25519PrivKey, err := crypto.GenerateAndEncodeICPIdentitySecretKey()
	if err != nil {
		return nil, err
	}

	// Write the ICP identity private key to the secrets manager storage
	if setErr := secretsManager.SetSecret(
		secrets.ICPIdentityKey,
		ed25519PrivKey,
	); setErr != nil {
		return nil, setErr
	}

	return ed25519PubKey, nil
}

func EncryptICPIdentityKey(secretsManager secrets.SecretsManager, secretsPass string) error {
	if secretsManager.HasSecret(secrets.SecureFlag + secrets.ICPIdentityKey) {
		secureFlag, err := secretsManager.GetSecret(
			secrets.SecureFlag + secrets.ICPIdentityKey,
		)
		if err != nil {
			return err
		}
		if string(secureFlag) == secrets.SecureTrue {
			return fmt.Errorf(`secrets "%s" has been already encrypted`, secrets.ICPIdentityKey)
		}
	}

	if secretsManager.HasSecret(secrets.ICPIdentityKey) {
		// encrypt ed25519 key for ICP identity
		icPrivKey, err := secretsManager.GetSecret(secrets.ICPIdentityKey)
		if err != nil {
			return err
		}
		encryptedKey, err := crypto.CFBEncrypt(string(icPrivKey), secretsPass)
		if err != nil {
			return err
		}
		// Write the ICP identity private key to the secrets manager storage
		if setErr := secretsManager.SetSecret(
			secrets.ICPIdentityKey,
			[]byte(encryptedKey),
		); setErr != nil {
			return setErr
		}
		if setErr := secretsManager.SetSecret(
			secrets.SecureFlag+secrets.ICPIdentityKey,
			[]byte(secrets.SecureTrue),
		); setErr != nil {
			return setErr
		}
	} else {
		// generate ed25519 key for ICP identity
		_, ed25519PrivKey, err := crypto.GenerateAndEncodeICPIdentitySecretKey()
		if err != nil {
			return err
		}
		encryptedKey, err := crypto.CFBEncrypt(string(ed25519PrivKey), secretsPass)
		if err != nil {
			return err
		}

		// Write the ICP identity private key to the secrets manager storage
		if setErr := secretsManager.SetSecret(
			secrets.ICPIdentityKey,
			[]byte(encryptedKey),
		); setErr != nil {
			return setErr
		}
		if setErr := secretsManager.SetSecret(
			secrets.SecureFlag+secrets.ICPIdentityKey,
			[]byte(secrets.SecureTrue),
		); setErr != nil {
			return setErr
		}
	}
	return nil
}
