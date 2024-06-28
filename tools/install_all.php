<?php

$cwd = getcwd();
chdir(__DIR__."/..");
passthru("bash tools/install_all.bash");
chdir($cwd);


