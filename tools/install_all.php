<?php

$cwd = getcwd();
chdir(__DIR__."/..");
passthru("bash tools/install_all.bash");
chdir($cwd);


@mkdir(__DIR__."/../keys",0777,true);
@mkdir(__DIR__."/../working_data/sessions",0777,true);
@mkdir(__DIR__."/../html/recieved_images",0777,true);

echo "\n";
echo "Make sure your keys are in keys folder \n";
echo "Make sure working_data/sessions is writable \n";
echo "Make sure html/recieved_images is writable \n";
echo "Start apache (or whatever) with html folder as document root on port 8000 ! \n";
echo "Start goserver_reverse_proxy in it's directory \n";




