<?php

$FILES_SEEN = [];


function process_json($json){
    global $FILES_SEEN;

    if(!isset($FILES_SEEN[$json->Path])){
        $FILES_SEEN[$json->Path] = true;
        if(!isset($json->Filedata) || $json->Filedata==""){
            $json->Filedata = @file_get_contents($json->Path);
        }
    }

    if(isset($json->Parent)){
        process_json($json->Parent);
    }
}

function flatten_json($json){
    $x = 0;
    $arr = [];
    while(true){
        $next = $json->Parent;
        $json->Parent = null;
        $arr[] = json_encode($json);
        
        if(!isset($next) || $next==null){
            break;
        }
        $x++;
        $json = $next;
    }
    return $arr;
}

function process_line($num, $line,$gz) {
    $json = json_decode($line);
    process_json($json);
    $jsons = flatten_json($json);
   // $line_out = json_encode($json)."\n";
    $line_out = implode("\t",$jsons)."\n";

    echo $line_out."\n\n";
    gzwrite ($gz,$line_out);
    //file_put_contents(__DIR__."/out/$num.html",$line_out."\n",FILE_APPEND);
}

$num = $argv[1];
$filename = __DIR__ . "/../working_data/ttd/$num.txt";

// Open the file for reading
$handle = fopen($filename, "r");

if ($handle) {
    $lineNumber = 0;
    
    $html = file_get_contents(__DIR__."/templates/html.html");
    file_put_contents(__DIR__."/out/$num.html",$html);
    $gz = gzopen(__DIR__."/out/$num.temp",'abw9');
    // Read the file line by line
    while (($line = fgets($handle)) !== false) {
        $lineNumber++;
        
        // Remove any carriage return characters
        $line = str_replace("\r", "", $line);
        
        // Trim any trailing newline characters
        $line = rtrim($line, "\n");
        
        // Process the line
        process_line($num, $line,$gz);
    }
    gzclose ($gz);

    // Close the file handle
    fclose($handle);



    $chunkSize = 6000; // NB! Must be divisible by 3, 
    // Open a.txt in append mode
    $aFile = fopen(__DIR__."/out/$num.html", 'ab');
    if ($aFile === false) {
        die("Could not open a.txt for writing.");
    }

    // Open b.bin in read mode
    $bFile = fopen(__DIR__."/out/$num.temp", 'rb');
    if ($bFile === false) {
        fclose($aFile);
        die("Could not open b.bin for reading.");
    }

    // Process the b.bin file in chunks
    while (!feof($bFile)) {
        // Read a chunk from b.bin
        $chunk = fread($bFile, $chunkSize);
        if ($chunk === false) {
            fclose($bFile);
            fclose($aFile);
            die("Error reading from b.bin.");
        }

        // Base64 encode the chunk
        $base64Chunk = base64_encode($chunk);

        // Write the base64 encoded chunk to a.txt
        $written = fwrite($aFile, $base64Chunk);
        if ($written === false) {
            fclose($bFile);
            fclose($aFile);
            die("Error writing to a.txt.");
        }
    }

    // Close both files
    fclose($bFile);
    fclose($aFile);

    unlink(__DIR__."/out/$num.temp");

} else {
    // Error opening the file
    echo "Error: Unable to open file '$filename'";
}