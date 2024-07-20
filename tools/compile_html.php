<?php
$cwd = getcwd();
chdir(__DIR__."/..");
passthru("bash tools/bash_scripts/compile_htmls.bash");
chdir($cwd);
