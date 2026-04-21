USE sts2_tracker;

SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE runs;
TRUNCATE TABLE run_cards;
TRUNCATE TABLE card_choice_options;
TRUNCATE TABLE map_point_card_choices;

SET FOREIGN_KEY_CHECKS = 1;