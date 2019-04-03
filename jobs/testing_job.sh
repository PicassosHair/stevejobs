#!/bin/bash

echo "This is a testing job that does not do anything"
/usr/src/app/jobs/log_slack info "Testing info msg."
/usr/src/app/jobs/log_slack success "Testing ok msg."
/usr/src/app/jobs/log_slack warning "Testing warning msg."
/usr/src/app/jobs/log_slack error "Testing bad msg."