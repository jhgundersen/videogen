#compdef videogen

_videogen() {
    local state

    _arguments \
        '1: :->model' \
        '*: :->args'

    case $state in
        model)
            _values 'model' \
                'seedance[ByteDance Seedance 2.0]' \
                'grok[xAI Grok Imagine Video]' \
                'sora[OpenAI Sora 2 Stable]'
            ;;
        args)
            case ${words[2]} in
                seedance)
                    _arguments \
                        '--duration[Duration in seconds]:duration:(5 10 15)' \
                        '--ratio[Aspect ratio]:ratio:(16:9 9:16 1:1 4:3 3:4 21:9)' \
                        '--image[Reference image URL]:url:' \
                        '--open[Open video after download]' \
                        '*:prompt:'
                    ;;
                grok)
                    _arguments \
                        '--duration[Duration in seconds]:duration:(10 15)' \
                        '--ratio[Aspect ratio]:ratio:(16:9 9:16 1:1 2:3 3:2)' \
                        '--image[Reference image URL]:url:' \
                        '--open[Open video after download]' \
                        '*:prompt:'
                    ;;
                sora)
                    _arguments \
                        '--duration[Duration in seconds]:duration:(10 15 25)' \
                        '--ratio[Aspect ratio]:ratio:(16:9 9:16)' \
                        '--variant[Model variant]:variant:(sora-2 sora-2-hd sora-2-pro)' \
                        '--image[Reference image URL]:url:' \
                        '--open[Open video after download]' \
                        '*:prompt:'
                    ;;
            esac
            ;;
    esac
}

_videogen
