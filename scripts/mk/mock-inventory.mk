##
# Rules to start / stop an inventory host mock
# server by using 'mockserver'
#
# See: https://github.com/mock-server/mockserver
# See: https://www.mock-server.com/
##
.PHONY: open-mockserver-ui
open-mockserver-ui:  ## Open the web console 
	xdg-open http://localhost:8010/mockserver/dashboard
