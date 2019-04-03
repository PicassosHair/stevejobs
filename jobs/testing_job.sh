#!/bin/bash

echo "This is a testing job that does not do anything"
/usr/src/app/jobs/log_slack.sh info "Testing info msg."
/usr/src/app/jobs/log_slack.sh success "Testing ok msg."
/usr/src/app/jobs/log_slack.sh warning "Testing warning msg."
/usr/src/app/jobs/log_slack.sh error "Testing bad msg."