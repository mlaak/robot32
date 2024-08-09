<?php

$cwd = getcwd();

chdir(__DIR__."/composer");
if(!file_exists("composer.phar")){
    passthru('curl "https://raw.githubusercontent.com/composer/getcomposer.org/76a7060ccb93902cd7576b67264ad91c8a2700e2/web/installer" > install_composer.php ');
    passthru('php  install_composer.php');
}

chdir(__DIR__."/..");
passthru("bash tools/bash_scripts/install_all.bash");
passthru("bash tools/bash_scripts/compile_htmls.bash");
chdir($cwd);


@mkdir(__DIR__."/../keys",0777,true);
@mkdir(__DIR__."/../working_data/sessions",0777,true);
@mkdir(__DIR__."/../working_data/ttd",0777,true);
@mkdir(__DIR__."/../collected_data/chats",0777,true);

@mkdir(__DIR__."/../html/recieved_images",0777,true);


echo "\n";
echo "Make sure your keys are in keys folder \n";
echo "Make sure working_data/sessions is writable \n";
echo "Make sure html/recieved_images is writable \n";
echo "Start apache (or whatever) with html folder as document root on port 8000 ! \n";
echo "Start goserver_reverse_proxy in it's directory \n";




