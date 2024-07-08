<?php
$cwd = getcwd();
chdir(__DIR__."/..");
passthru("bash tools/compile_htmls.bash");
chdir($cwd);
