# AutomatesFormalLanguages

Small Documentation:

- to run shell go to directory cmd/ and run command go build && ./cmd;
- commands:
    - modify
        - eps; flags: -n, name of automate; -h, help; usage: remove epsilon transitions;
        - det; flags: -n, name of automate, -h, help, usage: determine automate;
        - full; flags: -n, name of automate; -h, help; usage: make automate full;
        - min; flags: -n, name of automate; -h, help; usage: make automate min;
    - use
        - show; flags: -n, name of automate; -h, help; usage: get PNG image of automate;
        - read; flags: -n, name of automate; -h, help; usage: check if automate can read the word;
    - create 
        - automate; flags: -n, name of automate; -r, regular expression; -h, help; usage: build automate by expression