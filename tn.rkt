#lang racket/base

(require racket/list 
         racket/system 
         racket/cmdline 
         racket/match)

(define app-dir (or (getenv "TN_DIR") (build-path (getenv "HOME") ".tn")))
(define notes-dir (build-path app-dir "notes"))
(define editor (find-executable-path (or (getenv "EDITOR") "vim")))

(define (key->filename key #:extension [extension ".md"])
  (string-append key extension))

(define (strip-extension filename)
  (let ([match-result (regexp-match #px"(.*)[.][^.]*$" filename)])
    (if match-result
      (second match-result)
      filename)))

(define (key->filepath key)
  (build-path notes-dir (key->filename key)))

(define (edit-note key)
  (void (system* editor (key->filepath key))))

(define (list-notes)
  (for ([key (in-list (map strip-extension (list-dir-by-date-desc notes-dir)))])
    (displayln key)))

(define (list-dir-by-date-desc dir-path)
  (struct path/age (path age))
  (let ([path-ages
          (for/list [(p (in-list (directory-list dir-path)))]
            (path/age p (file-or-directory-modify-seconds (build-path dir-path p))))])
    (map
      path/age-path
      (sort
        path-ages
        (lambda (a b)
          (< (path/age-age a)
             (path/age-age b)))))))

(define (mkdirs-if-missing . dirs)
  (for ([dir (in-list dirs)])
    (unless (directory-exists? dir)
      (make-directory dir))))

(define (complete-notes args)
  (displayln args))

(define bash-completion-script/string
#<<END
_tn_complete() {
    local cur prev opts

    COMPREPLY=()

    cur=${COMP_WORDS[COMP_CWORD]}
    prev=${COMP_WORDS[COMP_CWORD-1]}
    opts=""

    case ${cur} in
        *)
            opts=$(tn --list-all)
            ;;
    esac

    if [[ ! -z ${opts} ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
}

complete -F _tn_complete tn
END
)

(define (exit/message message status)
  (displayln message)
  (exit status))

(module+ main
  (mkdirs-if-missing app-dir notes-dir)

  (define do-list-all (make-parameter #f))
  (define do-complete-notes (make-parameter #f))
  (define do-completion-script (make-parameter #f))

  (define args
    (command-line #:program "tn"
                  #:once-each [("-l" "--list-all") "list all notes" (do-list-all #t)] 
                              [("--bash-completion-script") "echo a completion script" (do-completion-script #t)]
                              
                  #:args args
                  args))

  (cond 
    [(do-list-all) (list-notes)]
    [(do-completion-script) (displayln bash-completion-script/string)]
    [(equal? (length args) 1) (edit-note (first args))]
    [else 
      (exit/message "Error: please provide a <note name>" 1)])
)
