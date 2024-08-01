# Improved Chainsaw
- Not a chainsaw
- Not an improvment

This is a simple golang regex parser and matcher using Thompson’s construction. It only supports matching the full input, and will not match substrings. It supports:
- Character class escapes (`\w`, `\d`, `\s`) and their negations (`\W`, `\D`, `\S`) 
- Character classes (`[A-Z]`, `[abc]`)
- Wildcard symbols (`.`)
- Quantifers (`*`, `+`, `?`, `{n}`, `{n,}`, `{n,m}`)
- Groups (`(abc)`)
- Alternation (`a|b`, `a|b|c`)
- Character Literals (`a`, `b`, `c`)