<?php
if(!function_exists("str_contains")){ //pre PHP 8
    function str_contains($haystack,$needle){
        return strpos($haystack,$needle)!==false;
    }
}
function run($cmd){
    $output = [];
    exec($cmd,$output);
    return strtolower(implode("\n",$output));
}
function runl($cmd){
    return strlen(run($cmd));
}

$G = "\033[1;32m"; //green text
$R = "\033[1;31m"; //red text
$N = "\033[0m"; //neutral - resets text

chdir(__DIR__."/../..");

$s=run('php expert_illustrator.php content="Hello" raw="base64" prerendered=1');
$i = strlen($s);
if($i<=10000){
    echo "$R ILLUSTRATOR EXPERT TEST FAIL1 (prerendered):$N ";
    echo "expected output size>10000, got $i: $s";
} else {
    echo "$G pass illustrator expert test 1 (prerendered)$N\n";
}

$s=run('php expert_illustrator.php content="Hello" raw="base64" no_prerendered=1');
$i = strlen($s);
if($i<=10000){
    echo "$R ILLUSTRATOR EXPERT TEST FAIL2 (non prerendered):$N ";
    echo "expected output size>10000, got $i: $s";
} else {
    echo "$G pass illustrator expert test 2 (non prerendered)$N\n";
}


$s=run('php expert_illustrator.php content="Hello" raw="base64" mess_with_key=1');
$i = strlen($s);
if($i<=10000){
    echo "$R ILLUSTRATOR EXPERT TEST FAIL3 (fallback):$N ";
    echo "expected output size>10000, got $i: $s";
} else {
    echo "$G pass illustrator expert test 3 (fallback)$N\n";
}
