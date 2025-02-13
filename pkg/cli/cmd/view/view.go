/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package view

import (
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/dnote/dnote/pkg/cli/cmd/cat"
	"github.com/dnote/dnote/pkg/cli/cmd/ls"
	"github.com/dnote/dnote/pkg/cli/utils"
)

var example = `
 * View all books
 dnote view

 * List notes in a book
 dnote view javascript

 * View a particular note in a book
 dnote view javascript 0
 `

var nameOnly bool
var contentOnly bool

func preRun(cmd *cobra.Command, args []string) error {
	return nil
}

// NewCmd returns a new view command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view <book name?> <note index?>",
		Aliases: []string{"v"},
		Short:   "List books, notes or view a content",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	f := cmd.Flags()
	f.BoolVarP(&nameOnly, "name-only", "", false, "print book names only")
	f.BoolVarP(&contentOnly, "content-only", "", false, "print the note content only")

	return cmd
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			run := ls.NewRun(ctx, nameOnly)
			return run(cmd, args)
		} else if len(args) == 1 {
			if nameOnly {
				return errors.New("--name-only flag is only valid when viewing books")
			}

			if utils.IsNumber(args[0]) {
				run := cat.NewRun(ctx, contentOnly)
				return run(cmd, args)
			} else {
				run := ls.NewRun(ctx, false)
				return run(cmd, args)
			}
		} else {
			// Manejo de múltiples IDs
			for _, arg := range args {
				if !utils.IsNumber(arg) {
					return errors.New("all arguments must be numeric note IDs")
				}

				// Ejecutar cat.NewRun para cada ID individualmente
				run := cat.NewRun(ctx, contentOnly)
				err := run(cmd, []string{arg})
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}
