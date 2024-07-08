<?php

function runPrograms($program1,$cwd1, $program2, $cwd2) {
    $descriptorspec = array(
        0 => array("pipe", "r"),  // stdin
        1 => STDOUT,  // stdout
        2 => STDERR   // stderr
    );

    // Start the first program
    $process1 = proc_open($program1, $descriptorspec, $pipes1, $cwd1);

    // Start the second program
    $process2 = proc_open($program2, $descriptorspec, $pipes2, $cwd2);

    // Check if the processes started successfully
    if (is_resource($process1) && is_resource($process2)) {
        // Wait for both processes to finish
        $status1 = proc_get_status($process1);
        $status2 = proc_get_status($process2);

        while ($status1['running'] || $status2['running']) {
            if ($status1['running']) {
                $status1 = proc_get_status($process1);
            }
            if ($status2['running']) {
                $status2 = proc_get_status($process2);
            }
            usleep(100000); // Sleep for 100ms to avoid excessive CPU usage
        }

        // Close the processes
        proc_close($process1);
        proc_close($process2);

        echo "Both programs have finished executing.\n";
    } else {
        echo "Failed to start one or both programs.\n";
    }
}

// Example usage
//$program1 = getExecutableCommand('program1');
//$program2 = getExecutableCommand('program2');

//runPrograms($program1,"", $program2);

/*
function getExecutableCommand($programName) {
    if (strtoupper(substr(PHP_OS, 0, 3)) === 'WIN') {
        // Windows
        return $programName . '.exe';
    } else {
        // Linux and macOS
        return './' . $programName;
    }
}*/


runPrograms("bash devserver_8080.bash", __DIR__."/../goserver_reverse_proxy/",
            "php -S localhost:8000",    __DIR__."/../html/"
);























