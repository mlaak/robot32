<?php
    $gz = gzopen ( "test.gz", 'abw9' );
    gzwrite ( $gz, "test" );
    gzclose ( $gz );
?>