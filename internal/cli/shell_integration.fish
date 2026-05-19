function cdbm
    if test (count $argv) -eq 0
        command cdbm $argv
        return
    end

    set arg $argv[1]

    switch $arg
        case "add" "list" "remove" "edit" "help" "init" "" "-*"
            command cdbm $argv
        case '*'
            set output (command cdbm $argv)
            if string match -q "cd *" $output
                eval $output
            else
                echo "Invalid output from cdbm" >&2
                return 1
            end
    end
end
