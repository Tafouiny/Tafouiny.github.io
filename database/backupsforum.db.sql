BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "categories" (
	"id_category"	integer NOT NULL,
	"content"	text NOT NULL,
	PRIMARY KEY("id_category" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "users" (
	"id_user"	integer NOT NULL,
	"email"	text NOT NULL,
	"password"	text NOT NULL,
	"username"	char(60) NOT NULL,
	PRIMARY KEY("id_user" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "sessions" (
	"token"	text NOT NULL,
	"username"	text NOT NULL,
	"expiry"	timestamp NOT NULL,
	PRIMARY KEY("token")
);
CREATE TABLE IF NOT EXISTS "post_categories" (
	"uuid_post_categories"	text NOT NULL,
	"id_category"	integer NOT NULL,
	PRIMARY KEY("uuid_post_categories","id_category"),
	FOREIGN KEY("id_category") REFERENCES "categories"("id_category"),
	FOREIGN KEY("uuid_post_categories") REFERENCES "posts"("uuid_post_categories")
);
CREATE TABLE IF NOT EXISTS "comments" (
	"id_comment"	integer NOT NULL,
	"id_post"	integer NOT NULL,
	"id_user"	integer NOT NULL,
	"comment_content"	text NOT NULL,
	PRIMARY KEY("id_comment" AUTOINCREMENT),
	FOREIGN KEY("id_post") REFERENCES "posts"("id_post"),
	FOREIGN KEY("id_user") REFERENCES "users"("id_user")
);
CREATE TABLE IF NOT EXISTS "posts" (
	"id_post"	integer NOT NULL,
	"post_title"	text NOT NULL,
	"post_content"	text NOT NULL,
	"uuid_post_categories"	text NOT NULL,
	"id_user"	integer NOT NULL,
	"time"	timestamp NOT NULL,
	PRIMARY KEY("id_post" AUTOINCREMENT),
	FOREIGN KEY("uuid_post_categories") REFERENCES "categories"("id_category"),
	FOREIGN KEY("id_user") REFERENCES "users"("id_user")
);
CREATE TABLE IF NOT EXISTS "post_likes" (
	"id_like"	integer NOT NULL,
	"id_post"	integer NOT NULL,
	"id_user"	integer NOT NULL,
	"like_type"	integer NOT NULL,
	UNIQUE("id_post","id_user"),
	PRIMARY KEY("id_like" AUTOINCREMENT),
	FOREIGN KEY("id_post") REFERENCES "posts"("id_post"),
	FOREIGN KEY("id_user") REFERENCES "users"("id_user")
);
CREATE TABLE IF NOT EXISTS "comment_likes" (
	"id_like_dislike"	integer NOT NULL,
	"id_comment"	integer NOT NULL,
	"id_user"	integer NOT NULL,
	"like_type"	integer NOT NULL,
	UNIQUE("id_comment","id_user"),
	PRIMARY KEY("id_like_dislike" AUTOINCREMENT),
	FOREIGN KEY("id_comment") REFERENCES "comments"("id_comment"),
	FOREIGN KEY("id_user") REFERENCES "users"("id_user")
);
INSERT INTO "categories" VALUES (1,'Politcs');
INSERT INTO "categories" VALUES (2,'Rap');
INSERT INTO "categories" VALUES (3,'Productivity');
INSERT INTO "categories" VALUES (4,'Programming');
INSERT INTO "categories" VALUES (5,'Database');
INSERT INTO "categories" VALUES (6,'Golang');
INSERT INTO "categories" VALUES (7,'Lifestyle');
INSERT INTO "users" VALUES (24,'xolsp123@gmail.com','$2a$10$BSo38fgEEjkp/GH1UEvAM.lHk2MIigR22VLL9cL/P21YH.k1.eT7q','Xander');
INSERT INTO "users" VALUES (25,'magniang@01talent.com','$2a$10$HNRUXerbmp39GUFxXso8v.yKq/vprFOzIb5LpSr7E/XpdjYFu0yWq','Niangoos');
INSERT INTO "users" VALUES (26,'max@gmail.com','$2a$10$UjKHeOxD3ZDtB.nQiM04DOGYCNyP2GUffSbsrbz2cV393pOVVbgty','Max');
INSERT INTO "users" VALUES (27,'niangpmaxgatte@gmail.com','$2a$10$fOpnfyL7rqsovtkyiq9NlOrxcWoZfytq16pEd8sxPUPij0i2cVF0e','Magatte');
INSERT INTO "users" VALUES (28,'sow06205@gmail.com','$2a$10$3whTwEIeIQvsHmaQ9AmjiuYBfDM5gtuqPnAGKI7qDsF2.JfnjjveC','jeebril_sow');
INSERT INTO "users" VALUES (29,'oumarlamsn@gmail.com','$2a$10$BTRsDh0WN4BjDkjdq9g2h.WLp7sjLcA3fjtkHWBVtdUCFLMrve6x2','oulam');
INSERT INTO "users" VALUES (30,'lomalack@gmail.com','$2a$10$sckSXGDzK/PfnOFmXkqcHOGI9ypLCTxIwTOM8bfB3dI0tRKWb9b0G','Bastian');
INSERT INTO "users" VALUES (31,'aya@gmail.com','$2a$10$n81bSczi9jCxjpKQECdYRO5rOcd3KWWeKA7RIyKZmyDzdTQn4nDIy','Aya');
INSERT INTO "sessions" VALUES ('2bb8bc65-6bd1-4269-9a9e-25310c12d6d0','oulam','2023-08-31 23:37:20.9580487+02:00');
INSERT INTO "post_categories" VALUES ('ecf7f544-b23f-477c-8793-f05ecf021a70',7);
INSERT INTO "post_categories" VALUES ('4182a528-23b5-4998-9db6-074aefd129a1',1);
INSERT INTO "post_categories" VALUES ('cdffcaf3-bbe2-4ac9-b2ed-8d9300280f39',3);
INSERT INTO "post_categories" VALUES ('5350dc7c-66f5-44e1-ad1f-5f051323678e',3);
INSERT INTO "post_categories" VALUES ('45e18ab3-bcdb-4695-b370-3b958c1ec0fb',3);
INSERT INTO "post_categories" VALUES ('45e18ab3-bcdb-4695-b370-3b958c1ec0fb',4);
INSERT INTO "post_categories" VALUES ('45e18ab3-bcdb-4695-b370-3b958c1ec0fb',5);
INSERT INTO "post_categories" VALUES ('45e18ab3-bcdb-4695-b370-3b958c1ec0fb',6);
INSERT INTO "post_categories" VALUES ('d38506a3-bdcb-49a3-816c-02925f10e8e6',1);
INSERT INTO "post_categories" VALUES ('d38506a3-bdcb-49a3-816c-02925f10e8e6',2);
INSERT INTO "post_categories" VALUES ('d38506a3-bdcb-49a3-816c-02925f10e8e6',3);
INSERT INTO "post_categories" VALUES ('d38506a3-bdcb-49a3-816c-02925f10e8e6',4);
INSERT INTO "post_categories" VALUES ('d38506a3-bdcb-49a3-816c-02925f10e8e6',5);
INSERT INTO "post_categories" VALUES ('d38506a3-bdcb-49a3-816c-02925f10e8e6',6);
INSERT INTO "post_categories" VALUES ('d38506a3-bdcb-49a3-816c-02925f10e8e6',7);
INSERT INTO "post_categories" VALUES ('e37faa39-f881-4483-b821-0dee83c5dc71',3);
INSERT INTO "post_categories" VALUES ('e1b6a5da-4df2-47bf-aa81-368d72dbf041',4);
INSERT INTO "post_categories" VALUES ('e1b6a5da-4df2-47bf-aa81-368d72dbf041',6);
INSERT INTO "post_categories" VALUES ('489fdbdd-7961-461a-adb0-d05fe5b4abe8',4);
INSERT INTO "post_categories" VALUES ('edafe31d-e699-43fc-87a7-6e6c416112f3',2);
INSERT INTO "post_categories" VALUES ('edafe31d-e699-43fc-87a7-6e6c416112f3',3);
INSERT INTO "post_categories" VALUES ('70f6cbc6-49bf-47db-b357-0d4b5cd37706',2);
INSERT INTO "post_categories" VALUES ('adaedb03-ab0e-407c-aa1a-4ed54d89bf04',2);
INSERT INTO "post_categories" VALUES ('adaedb03-ab0e-407c-aa1a-4ed54d89bf04',3);
INSERT INTO "post_categories" VALUES ('f265c3a3-445c-40e3-b6f4-e8b6e1327e46',2);
INSERT INTO "post_categories" VALUES ('f265c3a3-445c-40e3-b6f4-e8b6e1327e46',7);
INSERT INTO "post_categories" VALUES ('3be84839-a332-49af-a7e7-63408ed19d3e',2);
INSERT INTO "post_categories" VALUES ('3be84839-a332-49af-a7e7-63408ed19d3e',4);
INSERT INTO "post_categories" VALUES ('fa5ff471-d5d3-4e4c-99cb-8cf8415ba2fd',2);
INSERT INTO "post_categories" VALUES ('fa5ff471-d5d3-4e4c-99cb-8cf8415ba2fd',3);
INSERT INTO "post_categories" VALUES ('1e1897a6-a52e-4441-8080-ac8e1b3aa784',4);
INSERT INTO "post_categories" VALUES ('1e1897a6-a52e-4441-8080-ac8e1b3aa784',6);
INSERT INTO "post_categories" VALUES ('a6d6f128-b806-43ed-a2f8-631ef00fdce6',4);
INSERT INTO "post_categories" VALUES ('a6d6f128-b806-43ed-a2f8-631ef00fdce6',5);
INSERT INTO "post_categories" VALUES ('209da606-b758-4205-86ee-74d54d28648f',3);
INSERT INTO "post_categories" VALUES ('209da606-b758-4205-86ee-74d54d28648f',7);
INSERT INTO "post_categories" VALUES ('d185d8f6-597b-4ddb-889a-3e27fa78134a',1);
INSERT INTO "post_categories" VALUES ('d185d8f6-597b-4ddb-889a-3e27fa78134a',2);
INSERT INTO "post_categories" VALUES ('d185d8f6-597b-4ddb-889a-3e27fa78134a',3);
INSERT INTO "post_categories" VALUES ('d185d8f6-597b-4ddb-889a-3e27fa78134a',4);
INSERT INTO "post_categories" VALUES ('d185d8f6-597b-4ddb-889a-3e27fa78134a',5);
INSERT INTO "post_categories" VALUES ('d185d8f6-597b-4ddb-889a-3e27fa78134a',6);
INSERT INTO "post_categories" VALUES ('d185d8f6-597b-4ddb-889a-3e27fa78134a',7);
INSERT INTO "comments" VALUES (1,1,26,'ouais il fait vraiment beau ce soir');
INSERT INTO "comments" VALUES (2,1,24,'Je suis désespéré wala');
INSERT INTO "comments" VALUES (3,2,28,'Nice');
INSERT INTO "comments" VALUES (4,2,28,'haha');
INSERT INTO "comments" VALUES (5,3,28,'hi');
INSERT INTO "comments" VALUES (6,2,28,'haha');
INSERT INTO "comments" VALUES (7,4,28,'nice');
INSERT INTO "comments" VALUES (8,0,26,'hello');
INSERT INTO "comments" VALUES (9,0,26,'hello
');
INSERT INTO "comments" VALUES (10,0,26,'jj');
INSERT INTO "comments" VALUES (11,0,26,'hh');
INSERT INTO "comments" VALUES (12,0,26,'Max');
INSERT INTO "comments" VALUES (13,0,26,'Max');
INSERT INTO "comments" VALUES (14,6,26,'Helllo test');
INSERT INTO "comments" VALUES (15,5,26,'O like you so naz');
INSERT INTO "comments" VALUES (16,5,30,'Va te faire');
INSERT INTO "comments" VALUES (17,8,26,'');
INSERT INTO "comments" VALUES (18,8,26,'');
INSERT INTO "comments" VALUES (19,8,26,'max');
INSERT INTO "comments" VALUES (20,8,26,'amssssss');
INSERT INTO "comments" VALUES (21,8,26,'llaa');
INSERT INTO "comments" VALUES (22,8,26,'zal');
INSERT INTO "comments" VALUES (23,8,26,'blabla');
INSERT INTO "comments" VALUES (24,8,26,'hr');
INSERT INTO "comments" VALUES (25,8,26,'ee');
INSERT INTO "comments" VALUES (26,8,26,'eeee');
INSERT INTO "comments" VALUES (27,8,26,'eeeeee');
INSERT INTO "comments" VALUES (28,8,26,'eeeeee');
INSERT INTO "comments" VALUES (29,8,26,'eeee');
INSERT INTO "comments" VALUES (30,8,26,'eeee');
INSERT INTO "comments" VALUES (31,8,26,'azaz');
INSERT INTO "comments" VALUES (32,8,26,'papa');
INSERT INTO "comments" VALUES (33,8,29,'kl,npo');
INSERT INTO "comments" VALUES (34,8,29,'');
INSERT INTO "comments" VALUES (35,11,24,'coucou');
INSERT INTO "comments" VALUES (36,10,24,'Va t''en');
INSERT INTO "comments" VALUES (37,7,24,'hh');
INSERT INTO "comments" VALUES (38,12,24,'C''est mélangé');
INSERT INTO "comments" VALUES (39,12,26,'Plus jamais');
INSERT INTO "comments" VALUES (40,19,26,'I just wanna feel your skin on mine Feel your eyes do the exploring');
INSERT INTO "comments" VALUES (41,18,26,'Love is gone ');
INSERT INTO "comments" VALUES (42,20,26,'Ça saute');
INSERT INTO "posts" VALUES (1,'Belle vue','Vous trouvez pas qu''il fait beau chers amis','ecf7f544-b23f-477c-8793-f05ecf021a70',26,'2023-08-24 21:08:14.039727+00:00');
INSERT INTO "posts" VALUES (8,'Golang or Python','Which programming language for backend ? ','e1b6a5da-4df2-47bf-aa81-368d72dbf041',24,'2023-08-30 14:48:31.720725845+00:00');
INSERT INTO "posts" VALUES (10,'Dadju - Confessions ','Tu entendras sûrement dire
Que je n''suis pas l''homme le plus réglo
Que si je t''aime, je serais trop possessif
Que trop souvent j''agis par ego
C''est peut-être vrai, je n''me suis jamais senti
Comme quelqu''un digne d''être appelé "héros"
À mi-chemin entre ce qu''on appelle être gentil
Et peut-être le dernier des salauds','edafe31d-e699-43fc-87a7-6e6c416112f3',26,'2023-08-31 21:50:35.729605+00:00');
INSERT INTO "posts" VALUES (11,'Dadju - Confessions','Tu entendras sûrement dire
Que je n''suis pas l''homme le plus réglo
Que si je t''aime, je serais trop possessif
Que trop souvent j''agis par ego
C''est peut-être vrai, je n''me suis jamais senti
Comme quelqu''un digne d''être appelé "héros"
À mi-chemin entre ce qu''on appelle être gentil
Et peut-être le dernier des salauds','70f6cbc6-49bf-47db-b357-0d4b5cd37706',28,'2023-08-31 21:51:43.790391+00:00');
INSERT INTO "posts" VALUES (12,'Aya - Chérie ','J''me promenais tranquille
Et puis j''ai croisé, lui
Il m''a fait la cour, oui
Garçon m''a fait la cour, oui
J''me mets en valeur que pour un putain d''mec (sheesh)
D''après mes ex, j''suis une putain d''meuf (yeah)
Y a pas d''issue d''secours (oh non)
C''est la seule solution (yeah)
Igo, pardon (pardon)
C''est dans tes bras que tu veux que j''m''endorme
Attention, c''est dangereux, attention','adaedb03-ab0e-407c-aa1a-4ed54d89bf04',26,'2023-09-01 15:48:37.136656+00:00');
INSERT INTO "posts" VALUES (13,'Plus jamais - Aya ','''t''ai donné mon cœur, j''le referai plus jamais
J''ai trop de rancœur, ça n''arrivera plus jamais
J''ai déjà donné, j''le referai plus jamais
Ouais, j''ai déjà donné, ça m''arrivera plus jamais
J''t''ai donné mon cœur, j''le referai plus jamais
J''ai trop de rancœur, ça n''arrivera plus jamais
J''ai déjà donné, j''le referai plus jamais
Ouais, j''ai déjà donné, ça m''arrivera plus jamais','f265c3a3-445c-40e3-b6f4-e8b6e1327e46',26,'2023-09-01 16:14:32.399695+00:00');
INSERT INTO "posts" VALUES (14,'Aya - Djadja','Breuuuuuh','f837b6fb-0325-4cd2-88f5-45007cb76ed6',26,'2023-09-01 16:15:20.307772+00:00');
INSERT INTO "posts" VALUES (15,'Aya - Doudou ','T''es mimi, dis-le moi, doudou
Prouve-le moi, doudou
Et ça, c''est quel comportement, doudou?
Tu me mens beaucoup
Ça, c''est quel comportement, doudou?
Tu me mens beaucoup
Là, c''est clair, il m''fait tomber
J''le vois, laisse tomber
C''est qu''il est frais, j''vois flou
(C''est qu''il est frais, j''vois flou)
J''aime quand t''es là, c''est tout (ouais)
J''ai dit, "oh" (oh, oh, oh)
J''ai juré qu''on est plus que potos (oh, oh)
J''ai dit, "oh" (oh, oh, oh)
J''ai juré qu''on est plus que potos (oh, oh)','3be84839-a332-49af-a7e7-63408ed19d3e',31,'2023-09-01 16:18:01.937396+00:00');
INSERT INTO "posts" VALUES (16,'Aya - Jolie Nana','Jolie nana recherche joli djo
Comment on fait, j''suis pas très mytho
Moi j''ai le truc, je sens les pipeaux, ouais (Aya Nakamura)
Les pipeaux, oh yeah, yeah
Tu ne vas pas me manquer (non), toi qui croyais me connaître
T''es rempli de charabias, tu n''sauras plus rien de moi
Je ne peux pas supporter, t''as osé me comparer
T''inquiète pas, j''vais tout niquer, c''est la putain de life
','fa5ff471-d5d3-4e4c-99cb-8cf8415ba2fd',31,'2023-09-01 16:34:56.311917+00:00');
INSERT INTO "posts" VALUES (19,'Calvin Harris - One Kiss ','Let me take the night, I love real easy
And I know that you''ll still wanna see me
On the Sunday morning, music real loud
Let me love you while the moon is still out
Something in you, lit up heaven in me
The feeling won''t let me sleep
''Cause I''m lost in the way you move, the way you feel','209da606-b758-4205-86ee-74d54d28648f',26,'2023-09-01 16:37:47.842129+00:00');
INSERT INTO "post_likes" VALUES (95,7,30,-1);
INSERT INTO "post_likes" VALUES (96,6,30,1);
INSERT INTO "post_likes" VALUES (97,5,30,1);
INSERT INTO "post_likes" VALUES (98,8,28,1);
INSERT INTO "post_likes" VALUES (99,5,28,1);
INSERT INTO "post_likes" VALUES (100,7,24,1);
INSERT INTO "post_likes" VALUES (101,6,24,1);
INSERT INTO "post_likes" VALUES (110,9,29,1);
INSERT INTO "post_likes" VALUES (111,9,26,1);
INSERT INTO "post_likes" VALUES (121,9,28,1);
INSERT INTO "post_likes" VALUES (145,11,28,1);
INSERT INTO "post_likes" VALUES (147,10,28,1);
INSERT INTO "post_likes" VALUES (149,1,24,1);
INSERT INTO "post_likes" VALUES (154,11,24,1);
INSERT INTO "post_likes" VALUES (164,15,31,1);
INSERT INTO "post_likes" VALUES (165,12,31,1);
INSERT INTO "post_likes" VALUES (168,19,26,1);
INSERT INTO "post_likes" VALUES (169,20,26,1);
INSERT INTO "comment_likes" VALUES (5,0,24,1);
INSERT INTO "comment_likes" VALUES (8,35,24,1);
INSERT INTO "comment_likes" VALUES (9,34,24,-1);
INSERT INTO "comment_likes" VALUES (10,33,24,-1);
INSERT INTO "comment_likes" VALUES (11,32,24,-1);
INSERT INTO "comment_likes" VALUES (22,37,24,-1);
INSERT INTO "comment_likes" VALUES (23,5,31,1);
INSERT INTO "comment_likes" VALUES (24,40,26,1);
COMMIT;
