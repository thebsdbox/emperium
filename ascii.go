package main

import (
	"fmt"

	"github.com/fatih/color"
)

const tie = `
             ._,.
           "..-..pf.
          -L   ..#'
        .+_L  ."]#
        ,'j' .+.j'                 -'.__..,.,p.
       _~ #..<..0.                 .J-.''..._f.
      .7..#_.. _f.                .....-..,'4'
      ;' ,#j.  T'      ..         ..J....,'.j'
     .' .."^.,-0.,,,,yMMMMM,.    ,-.J...+'.j@
    .'.'...' .yMMMMM0M@^='""g.. .'..J..".'.jH
    j' .'1'  q'^)@@#"^".'"='BNg_...,]_)'...0-
   .T ...I. j"    .'..+,_.'3#MMM0MggCBf....F.
   j/.+'.{..+       '^~'-^~~""""'"""?'"'''1'
   .... .y.}                  '.._-:'_...jf 
   g-.  .Lg'                 ..,..'-....,'.
  .'.   .Y^                  .....',].._f
  ......-f.                 .-,,.,.-:--&'
                            .'...'..'_J'
                            .~......'#'        Tie Fighter Manufacturing   
                            '..,,.,_]'           Security Systems
                            .L..'..''.
`

type securityLevel struct {
	locks []bool
}

func InitSecurity(count int) *securityLevel {
	fmt.Println("Security Status>", color.YellowString("Initialising..."))
	fmt.Println("Security Status>", color.GreenString("Enabled"))
	return &securityLevel{locks: make([]bool, count)}
}

func (s *securityLevel) lockStatus(lock int) {
	if s.locks[lock] {
		fmt.Println(fmt.Sprintf("Security Lock %d> ", lock), color.WhiteString("["), color.RedString("▮"), color.WhiteString("]"))
	} else {
		fmt.Println(fmt.Sprintf("Security Lock %d> ", lock), color.WhiteString("["), color.GreenString("▮"), color.WhiteString("]"))
	}
}

func (s *securityLevel) Lock(lock int) {
	if s.locks[lock] {
		s.locks[lock] = false
		s.lockStatus(lock)
	}
}

func (s *securityLevel) Unlock(lock int) {
	if !s.locks[lock] {
		s.locks[lock] = true
		s.lockStatus(lock)
	}
}

func (s *securityLevel) Status() {
	for x := range s.locks {
		s.lockStatus(x)
	}
}
