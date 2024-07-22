<?php
require __DIR__."/../../vendor/Robot32lib/ClassTree/ClassTree.php";
use Robot32lib\ClassTree\ClassTree;
chdir(__DIR__."/../..");

$G = "\033[1;32m"; //green text
$R = "\033[1;31m"; //red text
$N = "\033[0m"; //neutral - resets text

exec('php expert_classifier.php content="I want to build a robot"',$output);
$result = implode("\n",$output); $output = [];
if(strpos($result,ClassTree::ELECTRONICS)===false){
    echo "$R CLASSIFIER EXPERT TEST FAIL1:$N ";
    echo "expected to be classified as electronic ".ClassTree::ELECTRONICS." but got".$result."\n";
} else {
    echo "$G pass classifier test 1 $N\n";
}

exec('php expert_classifier.php content="Make a snake game in python"',$output);
$result = implode("\n",$output); $output = [];
if(strpos($result,ClassTree::PYTHON)===false){
    echo "$R CLASSIFIER EXPERT TEST FAIL2:$N ";
    echo "expected to be classified as python ".ClassTree::PYTHON." but got".$result."\n";
} else {
    echo "$G pass classifier test 2 $N\n";
}

exec('php expert_classifier.php content="Tell me something funny"',$output);
$result = implode("\n",$output); $output = [];
if(strpos($result,ClassTree::JOKE)===false){
    echo "$R CLASSIFIER EXPERT TEST FAIL3:$N ";
    echo "expected to be classified as a joke ".ClassTree::JOKE." but got".$result."\n";
} else {
    echo "$G pass classifier test 3 $N\n";
}


//test slight misspelling
exec('php expert_classifier.php content="What is micropyton"',$output);
$result = implode("\n",$output); $output = [];
if(strpos($result,ClassTree::MICROPYTHON_ELECTRONICS)===false){
    echo "$R CLASSIFIER EXPERT TEST FAIL4:$N ";
    echo "expected to be classified as Micropyton ".ClassTree::MICROPYTHON_ELECTRONICS." but got ".$result."\n";
}else {
    echo "$G pass classifier test 4 $N\n ";
}

