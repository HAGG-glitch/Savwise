INSERT INTO transaction_categories(id, name, applies_to) VALUES
(gen_random_uuid(), 'Income', 'income'),
(gen_random_uuid(), 'Food', 'expense'),
(gen_random_uuid(), 'Transport', 'expense'),
(gen_random_uuid(), 'Airtime/Data', 'expense'),
(gen_random_uuid(), 'Education', 'expense'),
(gen_random_uuid(), 'Entertainment', 'expense'),
(gen_random_uuid(), 'Business', 'expense'),
(gen_random_uuid(), 'Healthcare', 'expense'),
(gen_random_uuid(), 'Savings', 'expense'),
(gen_random_uuid(), 'Other', 'expense')
ON CONFLICT (name) DO NOTHING;
