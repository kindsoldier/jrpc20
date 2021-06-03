
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go .c
OBJEXT= none

sbin_PROGRAMS = pmapp
pmapp_SOURCES = pmapp.go

GOFLAGS = -ldflags="-s -w"

EXTRA_pmapp_SOURCES = \
	pmapp.go

EXTRA_pmapp_SOURCES += \
	pmconfig/pmconfig.go.in


pmapp$(EXEEXT): $(pmapp_SOURCES) $(EXTRA_pmapp_SOURCES)
	$(GO) build $(GOFLAGS) -o pmapp$(EXEEXT) $(pmapp_SOURCES)

EXTRA_DIST = \
	README.md \
	go.mod \
	go.sum

EXTRA_DIST += \
	samples/pmapp.yml

clean-local:
	rm -rf autom4te.cache

install-data-local:
	test -z $(DESTDIR)$(APP_LIBDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LIBDIR)
	test -z $(DESTDIR)$(APP_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(APP_CONFDIR)
	test -z $(DESTDIR)$(APP_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LOGDIR)
	test -z $(DESTDIR)$(APP_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(APP_RUNDIR)
#	test -z $(DESTDIR)$(APP_DATADIR) || $(MKDIR_P) $(DESTDIR)$(APP_DATADIR)
	test -z $(DESTDIR)$(APP_LIBDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LIBDIR)/pmpublic
	cp -R pmpublic/ $(DESTDIR)$(APP_LIBDIR)/pmpublic/
#EOF
