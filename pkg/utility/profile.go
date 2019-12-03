package utility

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

// GetProfiles : return profiles selected in .aws/credentials
func GetProfiles(credentialsPath string) (profiles []string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return profiles, err
	}

	fpath := strings.Replace(credentialsPath, "~", home, 1)

	f, err := os.Open(fpath)
	if err != nil {
		log.Println(err)
		return profiles, err
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, 4096)

	var p string
	var t []string
	var profileWithAssumeRole string

	for {
		l, e := reader.ReadString('\n')

		if strings.HasPrefix(l, "[") {
			// profile tag line
			t = strings.Split(p, "=")
			profileWithAssumeRole = t[0]
			if profileWithAssumeRole != "" {
				profiles = append(profiles, p)
			}

			p = l[1 : len(l)-2]
		}

		if strings.HasPrefix(l, "role_arn") {
			// role_arn line
			p = fmt.Sprintf("%s|%s", p, ConvNewline(l, ""))
		}

		if strings.HasPrefix(l, "mfa_serial") {
			// mfa_serial line
			p = fmt.Sprintf("%s|%s", p, ConvNewline(l, ""))
		}

		if strings.HasPrefix(l, "source_profile") {
			// source_profile
			p = fmt.Sprintf("%s|%s", p, ConvNewline(l, ""))
		}

		if e != nil {
			profileWithAssumeRole = t[0]
			if profileWithAssumeRole != "" {
				profiles = append(profiles, p)
			}
			break
		}
	}
	return
}

// FinderProfile : return profile selected in .aws/credentials
func FinderProfile(profiles []string) (profile string, err error) {
	idx, err := fuzzyfinder.FindMulti(
		profiles,
		func(i int) string {
			return fmt.Sprintf("%s",
				profiles[i],
			)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			p := profiles[i]

			f := strings.Split(p, "|")
			return strings.Join(f, "\n")
		}),
	)

	if err != nil {
		log.Fatal(err)
		return profile, err
	}

	for _, i := range idx {
		profile = profiles[i]
	}

	return profile, nil
}
