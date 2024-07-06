-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS wu_occupation (
	id int NOT NULL,
	occupation_code TEXT NOT NULL,
	occupation_name TEXT NOT NULL,
	position_level INT NOT NULL,
	last_updated TEXT NOT NULL
);

INSERT INTO wu_occupation (id,occupation_code,occupation_name,position_level,last_updated) VALUES
	 (1,'Airline/Maritime Employee','Airline/Maritime Employee',1,'2020-06-25 00:00:00.000'),
	 (2,'Art/Entertainment/Media/Sports','Art/Entertainment/Media/Sports Professional',1,'2020-06-25 00:00:00.000'),
	 (3,'Civil/Government Employee','Civil/Government Employee',1,'2020-06-25 00:00:00.000'),
	 (4,'Domestic Helper','Domestic Helper',0,'2020-06-25 00:00:00.000'),
	 (5,'Driver','Driver',0,'2020-06-25 00:00:00.000'),
	 (6,'Teacher/Educator','Teacher/Educator',1,'2020-06-25 00:00:00.000'),
	 (7,'Hotel/Restaurant/Leisure','Hotel/Restaurant/Leisure Services Employee',1,'2020-06-25 00:00:00.000'),
	 (8,'Housewife/Child Care','Housewife/Child Care',0,'2020-06-25 00:00:00.000'),
	 (9,'IT and Tech Professional','IT and Tech Professional',1,'2020-06-25 00:00:00.000'),
	 (10,'Laborer-Agriculture','Laborer-Agriculture',0,'2020-06-25 00:00:00.000'),
	 (11,'Laborer-Construction','Laborer-Construction',0,'2020-06-25 00:00:00.000'),
	 (12,'Laborer-Manufacturing','Laborer-Manufacturing',0,'2020-06-25 00:00:00.000'),
	 (13,'Laborer- Oil/Gas/Mining','Laborer-Oil/Gas/Mining/Forestry',0,'2020-06-25 00:00:00.000'),
	 (14,'Medical/Health Care','Medical and Health Care Professional',1,'2020-06-25 00:00:00.000'),
	 (15,'Non-profit/Volunteer','Non-profit Volunteer',0,'2020-06-25 00:00:00.000'),
	 (16,'Cosmetic/Personal Care','Cosmetic/Personal Care Services',1,'2020-06-25 00:00:00.000'),
	 (17,'Law Enforcement/Military','Law Enforcement/Military Professional',1,'2020-06-25 00:00:00.000'),
	 (18,'Office Professional','Office Professional',1,'2020-06-25 00:00:00.000'),
	 (19,'Prof Svc Practitioner','Professional Service Practitioner',1,'2020-06-25 00:00:00.000'),
	 (20,'Religious/Church Servant','Religious/Church Servant',0,'2020-06-25 00:00:00.000'),
	 (21,'Retail Sales','Retail Sales',1,'2020-06-25 00:00:00.000'),
	 (22,'Retired','Retired',0,'2020-06-25 00:00:00.000'),
	 (23,'Sales/Insurance/Real Estate','Sales/Insurance/Real Estate Professional',1,'2020-06-25 00:00:00.000'),
	 (24,'Science/Research Professional','Science/Research Professional',1,'2020-06-25 00:00:00.000'),
	 (25,'Security Guard','Security Guard',0,'2020-06-25 00:00:00.000'),
	 (26,'Self-Employed','Self-Employed',0,'2020-06-25 00:00:00.000'),
	 (27,'Skilled Trade/Specialist','Skilled Trade/Specialist',1,'2020-06-25 00:00:00.000'),
	 (28,'Student','Student',0,'2020-06-25 00:00:00.000'),
	 (29,'Unemployed','Unemployed',0,'2020-06-25 00:00:00.000');


CREATE TABLE IF NOT EXISTS wu_employment_position_level (
	id INT NOT NULL,
	code TEXT NOT NULL,
	position_level TEXT NOT NULL,
	last_updated TEXT NOT NULL
);

INSERT INTO wu_employment_position_level (id,code,position_level,last_updated) VALUES
	 (1,'Entry Level','Entry Level','2020-06-27 00:00:00.000'),
	 (2,'Mid-Level/Supervisory/Management','Mid-Level/Supervisory/Management','2020-06-27 00:00:00.000'),
	 (3,'Senior Level/Executive','Senior Level/Executive','2020-06-27 00:00:00.000'),
	 (4,'Owner','Owner','2020-06-27 00:00:00.000');

CREATE TABLE IF NOT EXISTS wu_purpose_of_transaction (
	id int not null,
	code text not null,
	purpose text not null,
	last_updated text not null
);

INSERT INTO wu_purpose_of_transaction (id,code,purpose,last_updated) VALUES
	 (1,'Family Support/Living Expenses','Family Support/Living Expenses','2020-06-25 00:00:00.000'),
	 (2,'Saving/Investments','Saving/Investments','2020-06-25 00:00:00.000'),
	 (3,'Gift','Gift','2020-06-25 00:00:00.000'),
	 (4,'Goods & Services payment','Goods & Services payment','2020-06-25 00:00:00.000'),
	 (5,'Travel expenses','Travel expenses','2020-06-25 00:00:00.000'),
	 (6,'Education/School Fee','Education/School Fee','2020-06-25 00:00:00.000'),
	 (7,'Rent/Mortgage','Rent/Mortgage','2020-06-25 00:00:00.000'),
	 (8,'Emergency/Medical Aid','Emergency/Medical Aid','2020-06-25 00:00:00.000'),
	 (9,'Charity/Aid Payment','Charity/Aid Payment','2020-06-25 00:00:00.000'),
	 (10,'Employee Payroll/Employee Expense','Employee Payroll/Employee Expense','2020-06-25 00:00:00.000'),
	 (11,'Prize or Lottery Fees/Taxes','Prize or Lottery Fees/Taxes','2020-06-25 00:00:00.000');

CREATE TABLE IF NOT EXISTS wu_relationship (
	id int not null,
	code text not null,
	relationship text not null,
	last_updated text not null
);

INSERT INTO wu_relationship (id,code,relationship,last_updated) VALUES
	 (1,'Family','Family','2020-06-27 00:00:00.000'),
	 (2,'Friend','Friend','2020-06-27 00:00:00.000'),
	 (3,'Trade/BusinesPartner','Trade/Business Partner','2020-06-27 00:00:00.000'),
	 (4,'Employee/Employer','Employee/Employer','2020-06-27 00:00:00.000'),
	 (5,'Donor/Receiver of Ch','Donor/Receiver of Charitable Funds','2020-06-27 00:00:00.000'),
	 (6,'Purchaser/Seller','Purchaser/Seller','2020-06-27 00:00:00.000');

CREATE TABLE IF NOT EXISTS wu_source_of_funds (
	id int not null,
	code text not null,
	source_fund text not null,
	last_updated text not null
);

INSERT INTO wu_source_of_funds (id,code,source_fund,last_updated) VALUES
	 (1,'Salary','Salary','2020-07-13 08:51:11.643'),
	 (2,'Savings','Savings','2020-07-13 08:51:11.690'),
	 (3,'Borrowed Funds/Loan','Borrowed Funds/Loan','2020-07-13 08:51:11.737'),
	 (4,'Pension/Government/Welfare','Pension/Government/Welfare','2020-07-13 08:51:11.767'),
	 (5,'Gift','Gift','2020-07-13 08:51:11.813'),
	 (6,'Inheritance','Inheritance','2020-07-13 08:51:11.830'),
	 (7,'Charitable Donations','Charitable Donations','2020-07-13 08:51:11.877'),
	 (8,'Cash Tips','Cash Tips','2020-07-13 08:51:12.190'),
	 (9,'Sale of Goods/Property/Services','Sale of Goods/Property/Services','2020-07-13 08:51:12.267'),
	 (10,'Investment Income','Investment Income','2020-07-13 08:51:12.313');

CREATE TABLE IF NOT EXISTS wu_id_types (
	description text not null,
	code text not null,
	bo_version text not null
);

INSERT INTO wu_id_types (description,code,bo_version) VALUES
	 ('24k Card(Domestic Padala)','24K','0x00000000000A30ED'),
	 ('DSWD 4PS','4PS','0x00000000000B25B4'),
	 ('Alien Certification of Registration/Immigrant Certificate of Registration','ACR','0x00000000000A9970'),
	 ('AFP ID','AFP','0x00000000000B25B3'),
	 ('Barangay Certification','BCL','0x00000000000A997A'),
	 ('Company ID (Issuing company is regulated or supervised by BSP, SEC or IC)','CID','0x00000000000A9977'),
	 ('Certification from the National Council for the Welfare of Disabled Persons','CNC','0x00000000000A9971'),
	 ('Driver License','DLC','0x00000000000A9666'),
	 ('DSWD Certification','DSW','0x00000000000A9972'),
	 ('GSIS e-Card','GSS','0x00000000000A9979'),
	 ('Home Development Mutual Fund (HDMF) ID','HDM','0x00000000000A9973'),
	 ('IBP (Integrated Bar of the Philippines) ID','IBP','0x00000000000A9974'),
	 ('NBI Clearance','NBC','0x0000000000000FBA'),
	 ('New Tin ID','NTI','0x00000000000A96F0'),
	 ('OFW ID','OFW','0x00000000000A9975'),
	 ('Overseas Workers Welfare Administration (OWWA) ID','OWA','0x00000000000A9976'),
	 ('Police Clearance','PCL','0x0000000000000FBC'),
	 ('Postal ID','PID','0x0000000000000FBD'),
	 ('PAG IBIG ID/ HOME DEV. MUTUAL FUND','PII','0x00000000000B25B5'),
	 ('PRC ID','PRC','0x0000000000000FBE'),
	 ('Passport','PSS','0x0000000000000FBF'),
	 ('PWD ID /NTL. COUNCIL FOR THE WELFARE FOR DISABLED PERSONS)','PWD','0x00000000000B25B6'),
	 ('SAC ID','SAC','0x00000000000B2567'),
	 ('Senior Citizen ID','SCI','0x0000000000000FC0'),
	 ('Student ID (validated for school year)','SID','0x00000000000A9978'),
	 ('Seaman Book','SMB','0x00000000000A965E'),
	 ('SSS ID','SSS','0x0000000000000FC3'),
	 ('UMID','UMI','0x00000000000B25B7'),
	 ('Voter ID','VID','0x00000000000A962B');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS wu_occupation;
DROP TABLE IF EXISTS wu_employment_position_level;
DROP TABLE IF EXISTS wu_purpose_of_transaction;
DROP TABLE IF EXISTS wu_relationship;
DROP TABLE IF EXISTS wu_id_types;
DROP TABLE IF EXISTS wu_source_of_funds;


