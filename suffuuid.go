package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "suffuuid SUFFIX",
	Short: "suffuuid generates UUIDs with static suffixes",
	Long: `A generator for identifiable UUIDs with a shared suffix to make
identification, manipulation or removal easier in database data.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a UUID suffix")
		}

		if len(args) > 1 {
			return errors.New("unknown arguments provided")
		}

		if !isValidHexstring(args[0]) {
			return fmt.Errorf("invalid hex string provided: %s", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		count, err := cmd.Flags().GetInt("count")
		maybeDie("parsing args", err)

		var suffix []byte

		if len(args[0])%2 == 0 {
			suffix, err = hex.DecodeString(args[0])
		} else {
			suffix, err = hex.DecodeString("0" + args[0])
		}

		maybeDie("parsing hex string", err)

		randomSegmentEnd := 16 - len(suffix)

		for i := 0; i < count; i++ {
			randomUuid, err := uuid.NewRandom()
			maybeDie("generating randomness", err)

			randomBinary, err := randomUuid.MarshalBinary()
			maybeDie("converting to binary", err)

			constructed := append(randomBinary[0:randomSegmentEnd], suffix...)
			newUuid, err := uuid.FromBytes(constructed)
			maybeDie("constructing uuid", err)

			fmt.Println(newUuid.String())
		}
	},
}

func main() {
	root.PersistentFlags().Int("count", 1, "the number of UUIDs to generate")
	root.Execute()
}

func isValidHexstring(s string) bool {
	match, err := regexp.MatchString("^[0-9a-f]+$", s)
	maybeDie("checking suffix validity", err)

	return match
}

func maybeDie(where string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while %s: %v\n", where, err)
		os.Exit(1)
	}
}
