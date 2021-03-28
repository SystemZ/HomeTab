INSERT INTO `devices` (`id`, `user_id`, `name`, `token`, `token_push`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 1, 'PhoneX', 'ea2f5fc2-8fc9-11eb-b2f9-2fb84f6e38f4', 'v609234pclq0mamocusx3a', '2021-03-28 13:30:17', '2021-03-28 13:30:17', NULL),
(2, 1, 'Zphone', '0bae9ada-8fcc-11eb-96a4-bf39709cff3e', 'wcqd3uqj8jmf7q1xwevfzx', '2021-03-28 13:30:17', '2021-03-28 13:30:17', NULL);


INSERT INTO `events` (`id`, `user_id`, `device_id`, `code`, `val_str`, `val_int`, `created_at`) VALUES
(1, 1, 1, 1, NULL, 69, '2021-03-28 13:39:41'),
(6, 1, 1, 2, NULL, NULL, '2021-03-28 13:39:41'),
(9, 1, 1, 1, NULL, 42, '2021-03-28 13:39:41'),
(10, 1, 1, 3, 'Never Gonna Give You Up', NULL, '2021-03-28 13:39:41'),
(11, 1, 1, 4, 'Rick Astley', NULL, '2021-03-28 13:39:41'),
(12, 1, 1, 5, NULL, NULL, '2021-03-28 13:39:41');
