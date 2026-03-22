package starter

// Default skeletons per language when DB omits starter_code or a key.
var DefaultByLanguage = map[string]string{
	"python": "def solve():\n    pass\n\nif __name__ == \"__main__\":\n    solve()\n",
	"cpp": `#include <bits/stdc++.h>
using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    return 0;
}
`,
	"c": `#include <stdio.h>

int main(void) {
    return 0;
}
`,
	"java": `public class Main {
    public static void main(String[] args) {
    }
}
`,
	"javascript": `function solve() {
}
solve();
`,
}

// Merge returns a full map: defaults first, then DB overrides.
func Merge(db map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range DefaultByLanguage {
		out[k] = v
	}
	for k, v := range db {
		if v != "" {
			out[k] = v
		}
	}
	return out
}
