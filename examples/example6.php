<?php
// cronic:
//   name: Example PHP5 Job
//   desc: Say hello every 26 seconds
//   cron: */26 * * * * *
//   cmd: docker run --rm -v .:/app php:5.6-cli php /app/$f
echo "Hello, from PHP " . PHP_VERSION . "\n";
sleep(13);
echo "Bye, from PHP " . PHP_VERSION . "\n";
?>

