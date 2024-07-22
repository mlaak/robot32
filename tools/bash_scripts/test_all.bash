(cd goserver_reverse_proxy/limits &&  /usr/local/go/bin/go test)
(cd goserver_reverse_proxy/middlesitter &&  /usr/local/go/bin/go test)
(cd goserver_reverse_proxy/ratelimiter &&  /usr/local/go/bin/go test)
(cd goserver_reverse_proxy/situation &&  /usr/local/go/bin/go test)
(cd goserver_reverse_proxy/translator &&  /usr/local/go/bin/go test)
(cd goserver_reverse_proxy/usersession &&  /usr/local/go/bin/go test)

(cd html/experts/classifier/tests/integration && php test_classifier.php)
(cd html/experts/general/tests/integration && php test_general.php)
(cd html/experts/illustrator/tests/integration && php test_illustrator.php)
(cd html/experts/remote_llms/tests/integration && php test_remote_llms.php)
(cd html/experts/robotics/tests/integration && php test_robotics.php)



