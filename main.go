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
			tokenEnv, err := cmd.Flags().GetString("vault-wrapping-token-env-var")
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
			wrappingToken := os.Getenv(tokenEnv)
			if wrappingToken == "" {
				log.Println("INFO: no initial token found, exiting")
			}

			config := vault.DefaultConfig()
			if addressEnv != "" {
				config.Address = os.Getenv(addressEnv)
			}
			api, err := vault.NewClient(config)
			if err != nil {
				log.Fatal(err)
			}
			api.SetToken(wrappingToken)
			resp, err := api.Logical().Write("sys/wrapping/unwrap", nil)
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
	rootCmd.Flags().StringP("vault-wrapping-token-env-var", "t", "VAULT_WRAPPING_TOKEN", "vault wrapping token environment variable")
	rootCmd.Execute()
}
