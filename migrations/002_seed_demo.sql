INSERT INTO topics (id, slug) VALUES
    ('11111111-1111-1111-1111-111111111111', 'programming-languages'),
    ('22222222-2222-2222-2222-222222222222', 'data-analytics'),
    ('33333333-3333-3333-3333-333333333333', 'ai'),
    ('44444444-4444-4444-4444-444444444444', 'sql-and-databases'),
    ('55555555-5555-5555-5555-555555555555', 'soft-skills')
ON CONFLICT (id) DO NOTHING;

INSERT INTO topic_localizations (topic_id, locale, name, description) VALUES
    ('11111111-1111-1111-1111-111111111111', 'ru', 'Языки программирования', 'Курсы по языкам программирования и автоматизации.'),
    ('11111111-1111-1111-1111-111111111111', 'en', 'Programming languages', 'Courses about programming languages and automation.'),
    ('11111111-1111-1111-1111-111111111111', 'kz', 'Bagdarlamalau tilderi', 'Bagdarlamalau jane avtomattandyru turaly kurstar.'),
    ('22222222-2222-2222-2222-222222222222', 'ru', 'Аналитика данных', 'Курсы по анализу данных и визуализации.'),
    ('22222222-2222-2222-2222-222222222222', 'en', 'Data analytics', 'Courses on data analysis and visualization.'),
    ('22222222-2222-2222-2222-222222222222', 'kz', 'Derek taldauy', 'Derekterdi taldau jane vizualdandyru kurstary.'),
    ('33333333-3333-3333-3333-333333333333', 'ru', 'AI', 'Курсы по искусственному интеллекту и его применению.'),
    ('33333333-3333-3333-3333-333333333333', 'en', 'AI', 'Courses about artificial intelligence and practical use cases.'),
    ('33333333-3333-3333-3333-333333333333', 'kz', 'AI', 'Zhasandy intellekt jane ony qoldanu turaly kurstar.'),
    ('44444444-4444-4444-4444-444444444444', 'ru', 'SQL и базы данных', 'Курсы по SQL, PostgreSQL и проектированию данных.'),
    ('44444444-4444-4444-4444-444444444444', 'en', 'SQL and databases', 'Courses on SQL, PostgreSQL and data design.'),
    ('44444444-4444-4444-4444-444444444444', 'kz', 'SQL jane derekter qory', 'SQL, PostgreSQL jane derekterdi zhobalau turaly kurstar.'),
    ('55555555-5555-5555-5555-555555555555', 'ru', 'Soft skills', 'Курсы по коммуникации, лидерству и работе в команде.'),
    ('55555555-5555-5555-5555-555555555555', 'en', 'Soft skills', 'Courses on communication, leadership and teamwork.'),
    ('55555555-5555-5555-5555-555555555555', 'kz', 'Soft skills', 'Kommunikaciya, leadership jane komandalq jumys kurstary.')
ON CONFLICT (topic_id, locale) DO NOTHING;

INSERT INTO tags (id, slug) VALUES
    ('66666666-6666-6666-6666-666666666666', 'qa'),
    ('77777777-7777-7777-7777-777777777777', 'linux'),
    ('88888888-8888-8888-8888-888888888888', 'python'),
    ('99999999-9999-9999-9999-999999999999', 'sql'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'communication')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tag_localizations (tag_id, locale, name) VALUES
    ('66666666-6666-6666-6666-666666666666', 'ru', 'qa'),
    ('66666666-6666-6666-6666-666666666666', 'en', 'qa'),
    ('66666666-6666-6666-6666-666666666666', 'kz', 'qa'),
    ('77777777-7777-7777-7777-777777777777', 'ru', 'linux'),
    ('77777777-7777-7777-7777-777777777777', 'en', 'linux'),
    ('77777777-7777-7777-7777-777777777777', 'kz', 'linux'),
    ('88888888-8888-8888-8888-888888888888', 'ru', 'python'),
    ('88888888-8888-8888-8888-888888888888', 'en', 'python'),
    ('88888888-8888-8888-8888-888888888888', 'kz', 'python'),
    ('99999999-9999-9999-9999-999999999999', 'ru', 'sql'),
    ('99999999-9999-9999-9999-999999999999', 'en', 'sql'),
    ('99999999-9999-9999-9999-999999999999', 'kz', 'sql'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ru', 'коммуникация'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'en', 'communication'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'kz', 'kommunikaciya')
ON CONFLICT (tag_id, locale) DO NOTHING;

INSERT INTO courses (
    id, slug, status, level, duration_category, expected_hours, rating, rating_count, students_count,
    lessons_count, has_certificate, cover_image_url, author_name, published_at
) VALUES
    ('00000000-0000-0000-0000-000000000101', 'linux-for-qa-pipelines', 'published', 'intermediate', 'quick', 4, 4.5, 174, 174, 12, true, '', 'Ilia V.', NOW() - INTERVAL '20 days'),
    ('00000000-0000-0000-0000-000000000102', 'qa-communication-and-soft-skills', 'published', 'beginner', 'focused', 6, 4.3, 132, 132, 10, false, '', 'Amina S.', NOW() - INTERVAL '15 days'),
    ('00000000-0000-0000-0000-000000000103', 'data-analysis-with-sql', 'published', 'beginner', 'deep', 12, 4.7, 265, 265, 24, true, '', 'Dana K.', NOW() - INTERVAL '10 days'),
    ('00000000-0000-0000-0000-000000000104', 'practical-python-for-automation', 'published', 'intermediate', 'focused', 8, 4.8, 412, 412, 18, true, '', 'Murat T.', NOW() - INTERVAL '8 days'),
    ('00000000-0000-0000-0000-000000000105', 'ai-tools-for-educators', 'published', 'advanced', 'quick', 5, 4.6, 95, 95, 9, true, '', 'Aigerim N.', NOW() - INTERVAL '5 days'),
    ('00000000-0000-0000-0000-000000000106', 'postgresql-for-analysts', 'published', 'intermediate', 'focused', 9, 4.4, 210, 210, 16, false, '', 'Ruslan B.', NOW() - INTERVAL '3 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO course_localizations (course_id, locale, title, subtitle, short_description, description, syllabus) VALUES
    ('00000000-0000-0000-0000-000000000101', 'ru', 'Linux для QA Pipelines', 'Практический трек по Linux', 'Терминальные привычки, которые делают тестирование и CI стабильнее.', 'Курс знакомит с Linux-командами, bash-скриптами и приёмами, которые нужны инженеру по качеству в ежедневной работе.', '[{"title":"Linux basics"},{"title":"Bash for QA"},{"title":"CI troubleshooting"}]'),
    ('00000000-0000-0000-0000-000000000101', 'en', 'Linux for QA Pipelines', 'Practical Linux track', 'Terminal habits that make testing and CI more stable.', 'The course covers Linux commands, bash scripts and daily practices that help QA engineers work faster and safer.', '[{"title":"Linux basics"},{"title":"Bash for QA"},{"title":"CI troubleshooting"}]'),
    ('00000000-0000-0000-0000-000000000101', 'kz', 'QA pipelines ushin Linux', 'Praktikalyq Linux track', 'Testing jane CI jumysyn turaqty etetin terminal adetteri.', 'Kurs Linux komandalaryn, bash skriptterin jane QA mamandaryna kerek kun deli praktikany qamtidy.', '[{"title":"Linux basics"},{"title":"Bash for QA"},{"title":"CI troubleshooting"}]'),

    ('00000000-0000-0000-0000-000000000102', 'ru', 'QA Communication and Soft Skills', 'Коммуникация для начинающих специалистов', 'Как уверенно общаться в команде, на ревью и в переписке.', 'Курс помогает выстроить коммуникацию с разработчиками, менеджерами и заказчиками и уменьшить количество конфликтов.', '[{"title":"Team communication"},{"title":"Feedback"},{"title":"Conflict handling"}]'),
    ('00000000-0000-0000-0000-000000000102', 'en', 'QA Communication and Soft Skills', 'Communication for junior specialists', 'How to communicate clearly in teams, reviews and async chats.', 'This course helps learners build healthy communication with developers, managers and stakeholders.', '[{"title":"Team communication"},{"title":"Feedback"},{"title":"Conflict handling"}]'),
    ('00000000-0000-0000-0000-000000000102', 'kz', 'QA communication and soft skills', 'Junior mamandarga arnalgan kommunikaciya', 'Komandada, review-da jane chat-ta anyq soileu dağdylary.', 'Kurs learnerlerge damer men managerlermen tiimdi kommunikaciya quruğa komektesedi.', '[{"title":"Team communication"},{"title":"Feedback"},{"title":"Conflict handling"}]'),

    ('00000000-0000-0000-0000-000000000103', 'ru', 'Анализ данных с SQL', 'SQL для принятия решений', 'Практика выборок, агрегаций и бизнес-отчётов.', 'Курс строится вокруг реальных аналитических задач: сегментация, воронки, когортный анализ и интерпретация результатов.', '[{"title":"SELECT and JOIN"},{"title":"Aggregations"},{"title":"Analytical cases"}]'),
    ('00000000-0000-0000-0000-000000000103', 'en', 'Data Analysis with SQL', 'SQL for decision making', 'Hands-on queries, aggregations and business reports.', 'The course is built around analytical tasks such as segmentation, funnels, cohorts and result interpretation.', '[{"title":"SELECT and JOIN"},{"title":"Aggregations"},{"title":"Analytical cases"}]'),
    ('00000000-0000-0000-0000-000000000103', 'kz', 'SQL arqyly derek taldau', 'Sheshim qabyldauga arnalgan SQL', 'Query, aggregation jane business-report praktikasy.', 'Kurs segmentaciya, funnel, cohort jane natizheni tusindiru sekildi naqty tasktar negizinde qurylgan.', '[{"title":"SELECT and JOIN"},{"title":"Aggregations"},{"title":"Analytical cases"}]'),

    ('00000000-0000-0000-0000-000000000104', 'ru', 'Практический Python для автоматизации', 'Python без лишней теории', 'Автоматизация рутинных задач, API и тестовых сценариев.', 'Курс помогает быстро перейти от базового синтаксиса к реальным инструментам автоматизации и написанию полезных скриптов.', '[{"title":"Python basics"},{"title":"API automation"},{"title":"Useful scripts"}]'),
    ('00000000-0000-0000-0000-000000000104', 'en', 'Practical Python for Automation', 'Python without extra theory', 'Automate routine tasks, APIs and test scenarios.', 'The course helps learners move from syntax basics to practical automation scripts and workflows.', '[{"title":"Python basics"},{"title":"API automation"},{"title":"Useful scripts"}]'),
    ('00000000-0000-0000-0000-000000000104', 'kz', 'Avtomattandyru ushin praktikalyq Python', 'Artyq teoriiasyz Python', 'Rutina tasktar, API jane test scenario-lardy avtomattandyru.', 'Kurs sintaksis negizinen paydaly avtomattandyru skriptterine jyljam otyga komektesedi.', '[{"title":"Python basics"},{"title":"API automation"},{"title":"Useful scripts"}]'),

    ('00000000-0000-0000-0000-000000000105', 'ru', 'AI tools for educators', 'Инструменты ИИ для преподавателя', 'Как применять ИИ в подготовке уроков и материалов.', 'Курс показывает, как использовать AI для создания интерактивных заданий, обратной связи и персонализации обучения.', '[{"title":"Prompting basics"},{"title":"AI lesson plans"},{"title":"Content personalization"}]'),
    ('00000000-0000-0000-0000-000000000105', 'en', 'AI Tools for Educators', 'AI toolkit for teachers', 'How to use AI when preparing lessons and course materials.', 'The course demonstrates AI-assisted lesson planning, feedback generation and personalization workflows.', '[{"title":"Prompting basics"},{"title":"AI lesson plans"},{"title":"Content personalization"}]'),
    ('00000000-0000-0000-0000-000000000105', 'kz', 'Oqytushylar ushin AI tools', 'Mugalimge arnalgan AI toolkit', 'Sabaq pen material dayyndauda AI qoldanu joldary.', 'Kurs AI komegimen lesson planning, feedback jane personalization processterin korsetedі.', '[{"title":"Prompting basics"},{"title":"AI lesson plans"},{"title":"Content personalization"}]'),

    ('00000000-0000-0000-0000-000000000106', 'ru', 'PostgreSQL для аналитиков', 'Postgres без страха', 'От таблиц и индексов до оптимизации запросов.', 'Курс помогает понять, как устроен PostgreSQL, как проектировать таблицы и ускорять реальные аналитические запросы.', '[{"title":"Tables and indexes"},{"title":"Query plans"},{"title":"Performance basics"}]'),
    ('00000000-0000-0000-0000-000000000106', 'en', 'PostgreSQL for Analysts', 'Postgres without fear', 'From tables and indexes to query optimization.', 'This course explains how PostgreSQL works, how to model tables and how to speed up analytical queries.', '[{"title":"Tables and indexes"},{"title":"Query plans"},{"title":"Performance basics"}]'),
    ('00000000-0000-0000-0000-000000000106', 'kz', 'Analitikter ushin PostgreSQL', 'Qoryqysyz Postgres', 'Kesteler men indexterden query optimization-ga deiin.', 'Kurs PostgreSQL ishin tusinuge, keste modeldeuge jane analytical query-lardy tezdetuge komektesedi.', '[{"title":"Tables and indexes"},{"title":"Query plans"},{"title":"Performance basics"}]')
ON CONFLICT (course_id, locale) DO NOTHING;

INSERT INTO course_topics (course_id, topic_id) VALUES
    ('00000000-0000-0000-0000-000000000101', '11111111-1111-1111-1111-111111111111'),
    ('00000000-0000-0000-0000-000000000102', '55555555-5555-5555-5555-555555555555'),
    ('00000000-0000-0000-0000-000000000103', '22222222-2222-2222-2222-222222222222'),
    ('00000000-0000-0000-0000-000000000103', '44444444-4444-4444-4444-444444444444'),
    ('00000000-0000-0000-0000-000000000104', '11111111-1111-1111-1111-111111111111'),
    ('00000000-0000-0000-0000-000000000105', '33333333-3333-3333-3333-333333333333'),
    ('00000000-0000-0000-0000-000000000105', '55555555-5555-5555-5555-555555555555'),
    ('00000000-0000-0000-0000-000000000106', '44444444-4444-4444-4444-444444444444'),
    ('00000000-0000-0000-0000-000000000106', '22222222-2222-2222-2222-222222222222')
ON CONFLICT (course_id, topic_id) DO NOTHING;

INSERT INTO course_tags (course_id, tag_id) VALUES
    ('00000000-0000-0000-0000-000000000101', '66666666-6666-6666-6666-666666666666'),
    ('00000000-0000-0000-0000-000000000101', '77777777-7777-7777-7777-777777777777'),
    ('00000000-0000-0000-0000-000000000102', '66666666-6666-6666-6666-666666666666'),
    ('00000000-0000-0000-0000-000000000102', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('00000000-0000-0000-0000-000000000103', '99999999-9999-9999-9999-999999999999'),
    ('00000000-0000-0000-0000-000000000104', '88888888-8888-8888-8888-888888888888'),
    ('00000000-0000-0000-0000-000000000105', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('00000000-0000-0000-0000-000000000106', '99999999-9999-9999-9999-999999999999')
ON CONFLICT (course_id, tag_id) DO NOTHING;
