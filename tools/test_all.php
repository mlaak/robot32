<?php
$cwd = getcwd();

chdir(__DIR__."/..");
passthru("bash tools/bash_scripts/test_all.bash");

chdir($cwd);