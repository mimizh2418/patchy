package main

import "patchy/diff"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	ops, err := diff.Diff("testfiles/a.txt", "testfiles/b.txt")
	if err != nil {
		panic(err)
	}
	diff.PrintDiff(ops)
}
