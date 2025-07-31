<?php
// cronic:
//   name: Example PHP5 Job
//   desc: Say hello every 1 minute
//   cron: 0 * * * * *
//   cmd: docker run --rm -v .:/app php:5.6-cli php /app/$f
echo "Hello, from PHP " . PHP_VERSION . "\n";
sleep(30);
echo "Bye, from PHP " . PHP_VERSION . "\n";
?>

