package main

import (
	"io"
	"log"
	"os"
	"text/template"

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

const fileTemplate = `VAULT_TOKEN={{ .Token }}
VAULT_ADDR={{ .VaultAddr }}
`

func main() {
	tpl, err := template.New("").Parse(fileTemplate)
	if err != nil {
		log.Fatal(err)
	}
	rootCmd := &cobra.Command{
		Use:   "vault-config-extractor",
		Short: "extract vault addr and a vault token from environment variables",
		Run: func(cmd *cobra.Command, _ []string) {

			out, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatal(err)
			}
			addressEnv, err := cmd.Flags().GetString("vault-addr-env-var")
			if err != nil {
				log.Fatal(err)
			}
			roleIDEnv, err := cmd.Flags().GetString("vault-role-id-env")
			if err != nil {
				log.Fatal(err)
			}
			secretIDEnv, err := cmd.Flags().GetString("vault-secret-id-env")
			if err != nil {
				log.Fatal(err)
			}
			var writer io.Writer
			switch out {
			case "/dev/stdout":
				fallthrough
			case "-":
				writer = cmd.OutOrStdout()
			default:
				if _, err := os.Stat(out); !os.IsNotExist(err) {
					log.Println("INFO: outfile exists, exiting")
					return
				}
				fd, err := os.Create(out)
				if err != nil {
					log.Fatal(err)
				}
				defer fd.Close()
				writer = fd
			}
			roleID := os.Getenv(roleIDEnv)
			if roleID == "" {
				log.Println("INFO: no approle role id provided, exiting")
			}
			log.Printf("INFO: using role-id %s", roleID)
			secretID := os.Getenv(secretIDEnv)
			if roleID == "" {
				log.Println("INFO: no approle secret id provided, exiting")
			}

			config := vault.DefaultConfig()
			if addressEnv != "" {
				config.Address = os.Getenv(addressEnv)
			}
			api, err := vault.NewClient(config)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := api.Logical().Write("auth/approle/login", map[string]interface{}{
				"role_id":   roleID,
				"secret_id": secretID,
			})
			if err != nil {
				log.Fatal(err)
			}
			err = tpl.Execute(writer, map[string]interface{}{
				"Token":     resp.Auth.ClientToken,
				"VaultAddr": api.Address(),
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	rootCmd.Flags().StringP("out", "o", "-", "write output to this file")
	rootCmd.Flags().StringP("vault-addr-env-var", "a", "", "vault server address environment variable")
	rootCmd.Flags().StringP("vault-role-id-env", "r", "VAULT_APPROLE_ID", "vault approle role name environment variable")
	rootCmd.Flags().StringP("vault-secret-id-env", "s", "VAULT_APPROLE_SECRET_ID", "vault approle secret id environment variable")
	rootCmd.Execute()
}
