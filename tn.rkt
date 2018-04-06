#lang racket

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

(module+ main
  (mkdirs-if-missing app-dir notes-dir)

  (define args
    (command-line #:program "tn"
                  #:args args
                  args))

  (match args
    [(list key) (edit-note key)]
    [_ (list-notes)])
  )
