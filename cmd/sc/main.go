package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/austindizzy/securitycenter-go"
)

func exit() {
	fmt.Println("Exiting...")
	fmt.Println()
}

func printHelp() {
	fmt.Println("sc :	a SecurityCenter command line interface, written in Go\n\nUsage: sc [options] [action]\n\nOptions")
	flag.PrintDefaults()
	fmt.Println("Actions")
	for _, a := range actions {
		fmt.Printf("\t%s\t- %s\n", a.name, a.desc)
	}
}

var (
	host     string
	username string
	password string
	token    string
	session  string
	proxy    string
	verbose  bool
	skipSSL  bool
	timeout  int
	throttle string

	args []string

	fields     string
	outputType string

	actions []action
)

type action struct {
	name   string
	desc   string
	action func(*sc.SC)
}

func init() {
	flag.StringVar(&host, "host", "localhost", "the SecurityCenter host to connect to")
	flag.StringVar(&username, "user", "", "the username to use for any authentication requests")
	flag.StringVar(&password, "pass", "", "the password to use for any authentication requests")
	flag.StringVar(&token, "token", "", "the token to use for authenticated requests")
	flag.StringVar(&session, "session", "", "the session to use for authenticated requests")
	flag.StringVar(&fields, "fields", "", "comma-separted value of fields to include for applicable request types")
	flag.StringVar(&outputType, "type", "rich", "default cli output type (i.e. 'rich' for tables and graphic output, 'csv' for comma-separated values,'json' for raw JSON API responses) where supported")
	flag.StringVar(&proxy, "proxy", "", "the proxy to send all client communications over, supports HTTP/HTTPS/SOCKS")
	flag.BoolVar(&verbose, "verbose", false, "turn on verbose logging mode")
	flag.BoolVar(&skipSSL, "skipSSL", false, "disable SSL verification with SC instance requests")
	flag.IntVar(&timeout, "timeout", -1, "default number of seconds to wait before abandoning a request")
	flag.StringVar(&throttle, "throttle", "0s", "default minimum time duration to wait between making requests")

	flag.Lookup("timeout").DefValue = "90"

	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		exit()
		os.Exit(1)
	}()

	actions = []action{
		{"auth", "authenticate with the SecurityCenter instance, and receive a token and session", doAuth},
		{"whoami", "print user details for currently authenticated user (if one exists)", doWhoAmI},
		{"user", "print a specific user's information from the SecurityCenter instance", doUser},
		{"users", "print a list of users in the SecurityCenter instance", doUsers},
		{"status", "gets a collection of status information, including license information", doStatus},
	}
}

func do(s *sc.SC, action string) error {
	for _, a := range actions {
		if a.name == action {
			a.action(s)
			return nil
		}
	}
	return fmt.Errorf("the action \"%s\" does not exist", action)
}

func main() {
	flag.Parse()

	sc.Verbose = verbose
	sc.SkipSSLVerify = skipSSL
	if timeout > -1 {
		sc.TimeoutDuration = timeout
	}
	if len(proxy) > 1 {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "the proxy URL \"%s\" is invalid", proxy)
			os.Exit(1)
		}
		sc.Proxy = proxyURL
	}
	args = flag.Args()

	if len(args) == 0 {
		printHelp()
		return
	}

	s := sc.NewSC(host)
	if len(session) > 0 && len(token) > 0 {
		s = s.WithAuth(session, token)
	}

	err := do(s, args[0])
	if err != nil {
		panic(err)
	}
}
