AC_INIT(pmapp,0.1)
AM_INIT_AUTOMAKE
AC_PREFIX_DEFAULT(/usr/local)

PACKAGE=pmapp


AC_CHECK_PROG(HAVE_GO, go, true, false, /bin /usr/local/bin /usr/bin)
if test "x$HAVE_GO" = "xfalse"; then
    AC_MSG_ERROR([Requested program go not found])
fi
AC_PATH_PROG([go],[go])
AC_PATH_PROG([GO],[go])
dnl AC_SUBST(go, go)

AC_CHECK_PROG(HAVE_CP, cp, true, false, /bin /usr/local/bin /usr/bin)
if test "x$HAVE_CP" = "xfalse"; then
    AC_MSG_ERROR([Requested program cp not found])
fi
AC_PATH_PROG([CP],[cp])

AC_PROG_INSTALL

AC_CANONICAL_HOST

case $host_os in
    *freebsd* )
    AC_SUBST(ROOT_GROUP, "wheel")
    AM_CONDITIONAL(FREEBSD_OS, true)
    AM_CONDITIONAL(LINUX_OS, false)
    OSNAME=freebsd
    ROOT_GROUP=wheel
        ;;
    *linux* )
    AC_SUBST(ROOT_GROUP, "root")
    AM_CONDITIONAL(FREEBSD_OS, false)
    AM_CONDITIONAL(LINUX_OS, true)
    OSNAME=linux
    ROOT_GROUP=root
        ;;
esac

AM_CONDITIONAL(SYSTEMD, false)
if test -d /lib/systemd/system; then
    AM_CONDITIONAL(SYSTEMD, true)
fi

AC_ARG_ENABLE([devel-mode],
    AS_HELP_STRING([--enable-devel-mode], [Enable developmend mode]))


AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRCDIR=`pwd`
    enable_devel_mode=yes
])

test "x$prefix" == "xNONE" && prefix=$ac_default_prefix
PREFIX=$prefix

dnl --------------------------------------------------------------------------------------
AC_ARG_WITH(confdir,
    AS_HELP_STRING([--with-confdir=PATH],[set configuration dir to PATH (default: "${ac_default_prefix}"/etc/${PACKAGE})]),
    [ if test ! -z "$with_confdir" ; then
        case $with_confdir in
            /*)
                APP_CONFDIR="$with_confdir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an absolute path to --with-confdir=PATH)
                ;;
        esac
    else
        APP_CONFDIR="$prefix/etc/${PACKAGE}"
    fi ],
    [
    APP_CONFDIR="$prefix/etc/${PACKAGE}"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_CONFDIR="${SRCDIR}/"
    sysconfdir="${SRCDIR}/"
], [
    test "x$APP_CONFDIR" == "x/usr/etc/${PACKAGE}" && APP_CONFDIR="/etc/${PACKAGE}"
    test "x$prefix" == "x/usr" && sysconfdir="/etc"
])

AC_MSG_NOTICE(app_confdir set as ${APP_CONFDIR})

AC_DEFINE_UNQUOTED(APP_CONFDIR, "$APP_CONFDIR", [location of configuration files for ${PACKAGE}])
AC_SUBST(APP_CONFDIR, "$APP_CONFDIR")

AC_DEFINE_UNQUOTED(app_confdir, "$APP_CONFDIR", [location of configuration files for ${PACKAGE}])
AC_SUBST(app_confdir, "$APP_CONFDIR")

dnl --------------------------------------------------------------------------------------
AC_ARG_WITH(logdir,
    AS_HELP_STRING([--with-logdir=PATH],[set file path for source logdir (default: /var/log/${PACKAGE}/${PACKAGE}.log)]),
    [ if test ! -z "$with_logdir" ; then
        case $with_logdir in
            /*)
                APP_LOGDIR="$with_logdir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an absolute path to --with-logdir=PATH)
                ;;
        esac
    else
        APP_LOGDIR="/var/log/${PACKAGE}"
    fi ],
    [
    APP_LOGDIR="/var/log/${PACKAGE}"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_LOGDIR="${SRCDIR}/log"
])

AC_MSG_NOTICE(app_logdir set as ${APP_LOGDIR})

AC_DEFINE_UNQUOTED(APP_LOGDIR, "$APP_LOGDIR", [location of ${PACKAGE} logdir])
AC_SUBST(APP_LOGDIR, "$APP_LOGDIR")

AC_DEFINE_UNQUOTED(app_logdir, "$APP_LOGDIR", [location of ${PACKAGE} logdir])
AC_SUBST(app_logdir, "$APP_LOGDIR")

dnl --------------------------------------------------------------------------------------
AC_ARG_WITH(rundir,
    AS_HELP_STRING([--with-rundir=PATH],[set file path for source rundir (default: /var/run/${PACKAGE})]),
    [ if test ! -z "$with_rundir" ; then
        case $with_rundir in
            /*)
                APP_RUNDIR="$with_rundir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an absolute path to --with-rundir=PATH)
                ;;
        esac
    else
        APP_RUNDIR="/var/run/${PACKAGE}"
    fi ],
    [
    APP_RUNDIR="/var/run/${PACKAGE}"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_RUNDIR="${SRCDIR}/run"
])

AC_MSG_NOTICE(app_rundir set as ${APP_RUNDIR})

AC_DEFINE_UNQUOTED(APP_RUNDIR, "$APP_RUNDIR", [location of pid file])
AC_SUBST(APP_RUNDIR, "$APP_RUNDIR")

AC_DEFINE_UNQUOTED(app_rundir, "$APP_RUNDIR", [location of pid file])
AC_SUBST(app_rundir, "$APP_RUNDIR")


dnl --------------------------------------------------------------------------------------

AC_ARG_WITH(libdir,
    AS_HELP_STRING([--with-libdir=PATH],[set file path for source libdir (default: ${PREFIX}/share/${PACKAGE})]),
    [ if test ! -z "$with_libdir" ; then
        case $with_libdir in
            /*)
                APP_LIBDIR="$with_libdir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an absolute path to --with-libdir=PATH)
                ;;
        esac
    else
        APP_LIBDIR="${PREFIX}/share/${PACKAGE}"
    fi ],
    [
    APP_LIBDIR="${PREFIX}/share/${PACKAGE}"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_LIBDIR="${SRCDIR}/"
])

AC_MSG_NOTICE(app_libdir set as ${APP_LIBDIR})

AC_DEFINE_UNQUOTED(APP_LIBDIR, "$APP_LIBDIR", [location of libs])
AC_SUBST(APP_LIBDIR, "$APP_LIBDIR")

AC_DEFINE_UNQUOTED(app_libdir, "$APP_LIBDIR", [location of libs])
AC_SUBST(app_libdir, "$APP_LIBDIR")

dnl --------------------------------------------------------------------------------------

case $host_os in
    *freebsd* )
    default_user="www"
    default_group="www"
        ;;
    *linux* )
    default_user="www-data"
    default_group="www-data"
        ;;
esac

AC_ARG_WITH(user,
    AS_HELP_STRING([--with-user=${PACKAGE}],[set executing user name]),
    [ if test ! -z "$with_user" ; then
        case $with_user in
            "")
                AC_MSG_ERROR(You must specify user name)
                ;;
            *)
                APP_USER="$with_user"
                ;;
        esac
    else
        APP_USER="$default_user"
    fi ],
    [ APP_USER="$default_user" ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_USER="`id -un`"
])

AC_MSG_NOTICE(app_user set as ${APP_USER})


AC_DEFINE_UNQUOTED(APP_USER, "$APP_USER", [effective user])
AC_SUBST(APP_USER, "$APP_USER")

AC_DEFINE_UNQUOTED(app_user, "$APP_USER", [effective user])
AC_SUBST(app_user, "$APP_USER")

dnl --------------------------------------------------------------------------------------

AC_ARG_WITH(group,
    AS_HELP_STRING([--with-group=${PACKAGE}],[set executing group name]),
    [ if test ! -z "$with_group" ; then
        case $with_group in
            "")
                AC_MSG_ERROR(You must specify group name)
                ;;
            *)
                APP_GROUP="$with_group"
                ;;
        esac
    else
        APP_GROUP="$default_group"
    fi ],
    [ APP_GROUP="$default_group" ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_GROUP="`id -gn`"
])

AC_MSG_NOTICE(app_group set as ${APP_GROUP})

AC_DEFINE_UNQUOTED(APP_GROUP, "$APP_GROUP", [effective group id])
AC_SUBST(APP_GROUP, "$APP_GROUP")

AC_DEFINE_UNQUOTED(app_group, "$APP_GROUP", [effective group id])
AC_SUBST(app_group, "$APP_GROUP")

dnl --------------------------------------------------------------------------------------

default_datadir="/var/db/$PACKAGE"

AC_ARG_WITH(datadir,
    AS_HELP_STRING([--with-datadir=PATH],[set data directory (default: $default_datadir)]),
    [ if test ! -z "$with_datadir" ; then
        case $with_datadir in
            /*)
                APP_DATADIR="$with_datadir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an absolute path to --with-datadir=PATH)
                ;;
        esac
    else
        APP_DATADIR="$default_datadir"
    fi ],
    [
    APP_DATADIR="$default_datadir"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_DATADIR="${SRCDIR}"
])

AC_MSG_RESULT(datadir set as ${APP_DATADIR})


AC_DEFINE_UNQUOTED(APP_DATADIR, "$APP_DATADIR", [location of database dir])
AC_SUBST(APP_DATADIR, "$APP_DATADIR")

AC_DEFINE_UNQUOTED(app_datadir, "$APP_DATADIR", [location of database dir])
AC_SUBST(app_datadir, "$APP_DATADIR")

dnl --------------------------------------------------------------------------------------

default_cachedir="/var/cache/$PACKAGE"

AC_ARG_WITH(cachedir,
    AS_HELP_STRING([--with-cachedir=PATH],[set data directory (default: $default_cachedir)]),
    [ if test ! -z "$with_cachedir" ; then
        case $with_cachedir in
            /*)
                APP_CACHEDIR="$with_cachedir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an absolute path to --with-cachedir=PATH)
                ;;
        esac
    else
        APP_CACHEDIR="$default_cachedir"
    fi ],
    [
    APP_CACHEDIR="$default_cachedir"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_CACHEDIR="${SRCDIR}/cache"
])

AC_MSG_RESULT(cachedir set as ${APP_CACHEDIR})


AC_DEFINE_UNQUOTED(APP_CACHEDIR, "$APP_CACHEDIR", [location of cache dir])
AC_SUBST(APP_CACHEDIR, "$APP_CACHEDIR")

AC_DEFINE_UNQUOTED(app_cachedir, "$APP_CACHEDIR", [location of cache dir])
AC_SUBST(app_cachedir, "$APP_CACHEDIR")

APP_DAEMONIZE="true"
AS_IF([test "x$enable_devel_mode" = "xyes"], [
    APP_DAEMONIZE="false"
])

AC_SUBST(app_daemonize, "$APP_DAEMONIZE")

AC_SUBST(app_name, "$PACKAGE")
dnl AC_SUBST(app_port, 8080)

AC_OUTPUT([
Makefile
pmconfig/pmconfig.go
])
dnl EOF
