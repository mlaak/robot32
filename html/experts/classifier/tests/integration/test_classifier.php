<?php
require __DIR__."/../../vendor/Robot32lib/ClassTree/ClassTree.php";
use Robot32lib\ClassTree\ClassTree;

if(!function_exists("str_contains")){ //pre PHP 8
    function str_contains($haystack,$needle){
        return strpos($haystack,$needle)!==false;
    }
}
function run($cmd){
    $output = [];
    exec($cmd,$output);
    return implode("\n",$output);
}
$G = "\033[1;32m"; //green text
$R = "\033[1;31m"; //red text
$N = "\033[0m"; //neutral - resets text


chdir(__DIR__."/../..");



$r=run('php expert_classifier.php content="I want to build a robot"');
if(strpos($r,ClassTree::ELECTRONICS)===false){
    echo "$R CLASSIFIER EXPERT TEST FAIL1(electronics):$N ";
    echo "expected to be classified as electronic ".ClassTree::ELECTRONICS." but got".$r."\n";
} else {
    echo "$G pass classifier test 1 (electronics) $N\n";
}

$r=run('php expert_classifier.php content="Make a snake game in python"');
if(!str_contains($r,ClassTree::PYTHON)){
    echo "$R CLASSIFIER EXPERT TEST FAIL2(python):$N ";
    echo "expected to be classified as python ".ClassTree::PYTHON." but got".$r."\n";
} else {
    echo "$G pass classifier test 2 (python) $N\n";
}

$r = run('php expert_classifier.php content="Tell me something funny"');
if(!str_contains($r,ClassTree::JOKE)){
    echo "$R CLASSIFIER EXPERT TEST FAIL3(joke):$N ";
    echo "expected to be classified as a joke ".ClassTree::JOKE." but got".$r."\n";
} else {
    echo "$G pass classifier test 3 (joke) $N\n";
}


//test slight misspelling
$r = run('php expert_classifier.php content="What is micropyton"');
if(!str_contains($r,ClassTree::MICROPYTHON_ELECTRONICS)){
    echo "$R CLASSIFIER EXPERT TEST FAIL4(slight mispelling of micropython): $N ";
    echo "expected to be classified as Micropyton ".ClassTree::MICROPYTHON_ELECTRONICS." but got ".$r."\n";
}else {
    echo "$G pass classifier test 4 (slight mispelling of micropython) $N\n ";
}

