_videogen() {
    local cur prev words cword
    _init_completion || return

    if [[ $cword -eq 1 ]]; then
        COMPREPLY=($(compgen -W "seedance grok sora" -- "$cur"))
        return
    fi

    case "${words[1]}" in
        seedance)
            case "$prev" in
                --duration) COMPREPLY=($(compgen -W "5 10 15" -- "$cur")) ;;
                --ratio)    COMPREPLY=($(compgen -W "16:9 9:16 1:1 4:3 3:4 21:9" -- "$cur")) ;;
                *)          COMPREPLY=($(compgen -W "--duration --ratio --image --open" -- "$cur")) ;;
            esac ;;
        grok)
            case "$prev" in
                --duration) COMPREPLY=($(compgen -W "10 15" -- "$cur")) ;;
                --ratio)    COMPREPLY=($(compgen -W "16:9 9:16 1:1 2:3 3:2" -- "$cur")) ;;
                *)          COMPREPLY=($(compgen -W "--duration --ratio --image --open" -- "$cur")) ;;
            esac ;;
        sora)
            case "$prev" in
                --duration) COMPREPLY=($(compgen -W "10 15 25" -- "$cur")) ;;
                --ratio)    COMPREPLY=($(compgen -W "16:9 9:16" -- "$cur")) ;;
                --variant)  COMPREPLY=($(compgen -W "sora-2 sora-2-hd sora-2-pro" -- "$cur")) ;;
                *)          COMPREPLY=($(compgen -W "--duration --ratio --variant --image --open" -- "$cur")) ;;
            esac ;;
    esac
}

complete -F _videogen videogen
