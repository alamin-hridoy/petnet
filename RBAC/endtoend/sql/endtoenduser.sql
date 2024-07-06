/* usermgmdb */
INSERT INTO "public"."organization_information" 
    ("id", "org_name", "contact_email", "contact_phone", "active", "created", "updated", "deleted") 
    VALUES 
    ('f5528d02-e6f8-40c6-9dd6-4531e96867ed', 'endtoend', 'endtoend@example.com', '', 't', '2021-05-06 11:09:34.89183+00', '2021-05-06 11:09:34.89183+00', NULL);

INSERT INTO "public"."user_account" 
    ("id", "org_id", "username", "password", "first_name", "last_name", "email", "email_verified", "invite_status", "invite_sender", "invite_code", "invite_expiry", "created", "updated", "deleted") 
    VALUES 
    ('817b2362-0aa0-4819-a323-2001f1984d3f', 'f5528d02-e6f8-40c6-9dd6-4531e96867ed', 'endtoend@example.com', 'K547X6dHM1XFUaTqPeDFG30Fz147NeCBdc24e5e669b8b9431c814ea7766dec50c25e5494b57a6841b623454d03779e92', 'endtoend', 'test', 'endtoend@example.com', 'f', 'Approved', '', ' ', '0001-01-01 00:00:00+00', '2021-05-06 11:09:34.922709+00', '2021-05-06 11:09:34.922709+00', NULL);

INSERT INTO "public"."service_account" 
    ("auth_type", "org_id", "environment", "client_name", "client_id", "create_user_id", "disable_user_id", "created", "disabled", "challenge") 
    VALUES 
    ('oauth', 'f5528d02-e6f8-40c6-9dd6-4531e96867ed', '', 'EndtoendTest', '59SUgTDANlU0MxvHeOOAYmO9','817b2362-0aa0-4819-a323-2001f1984d3f', '', '2021-05-06 11:09:35.08893+00', NULL, '');

INSERT INTO "public"."roles" 
    ("id", "org_id", "role_name", "description", "create_user_id", "created", "updated", "delete_user_id", "deleted") 
    VALUES 
    ('78ffaeeb-3305-4ed9-b883-4623cbeed2f6', 'f5528d02-e6f8-40c6-9dd6-4531e96867ed', 'Endtoend', 'Endtoend test', '817b2362-0aa0-4819-a323-2001f1984d3f', '2021-05-06 11:09:34.955934+00', '2021-05-06 11:09:34.955934+00', NULL, NULL);

INSERT INTO "public"."permissions" 
    ("id", "org_id", "service_permission_id", "permission_name", "description", "create_user_id", "created", "updated", "delete_user_id", "deleted") 
    VALUES 
    ('ac1fa1c8-bf00-46bf-b99f-2834e6a42fa3', 'f5528d02-e6f8-40c6-9dd6-4531e96867ed', '6cbd6871-ea0f-48b2-b00c-9bb13b0c1623', 'Account Endtoend', 'Create System Accounts', '817b2362-0aa0-4819-a323-2001f1984d3f', '2021-05-06 11:09:34.954102+00', '2021-05-06 11:09:34.954102+00', NULL, NULL), 
    ('fa2a47cd-109a-496a-b5ee-f78f4f927c56', 'f5528d02-e6f8-40c6-9dd6-4531e96867ed', 'b2c04379-6329-4ce7-998c-2d6b22e2e897', 'Permission Endtoend', 'Create System Permissions', '817b2362-0aa0-4819-a323-2001f1984d3f', '2021-05-06 11:09:34.952011+00', '2021-05-06 11:09:34.952011+00', NULL, NULL);

/* hydradb */
INSERT INTO "public"."hydra_client" 
    ("id", "client_name", "client_secret", "redirect_uris", "grant_types", "response_types", "scope", "owner", "policy_uri", "tos_uri", "client_uri", "logo_uri", "contacts", "client_secret_expires_at", "sector_identifier_uri", "jwks", "jwks_uri", "request_uris", "token_endpoint_auth_method", "request_object_signing_alg", "userinfo_signed_response_alg", "subject_type", "allowed_cors_origins", "pk", "audience", "created_at", "updated_at", "frontchannel_logout_uri", "frontchannel_logout_session_required", "post_logout_redirect_uris", "backchannel_logout_uri", "backchannel_logout_session_required", "metadata", "token_endpoint_auth_signing_alg") 
    VALUES 
    ('59SUgTDANlU0MxvHeOOAYmO9', '', '$2a$10$f4s4GQRfAClCNJLjDYM4h.oDmkGVKYa3tNMYofAj.rh4ndqfvMl9C', '', 'client_credentials', 'code|id_token', 'offline_access offline openid', '', '', '', '', '', '', 0, '', '{}', '', '', 'client_secret_basic', '', 'none', 'public', '', 1, '', '2021-05-06 11:09:35', '2021-05-06 11:09:35', '', 'f', '', '', 'f', '{}', '');

/* ketodb */
INSERT INTO "public"."rego_data" 
    ("id", "collection", "pkey", "document") 
    VALUES 
    (1, '/store/ory/exact/policies', 'fa2a47cd-109a-496a-b5ee-f78f4f927c56', '{"id":"fa2a47cd-109a-496a-b5ee-f78f4f927c56","description":"Create System Permissions","subjects":["817b2362-0aa0-4819-a323-2001f1984d3f","78ffaeeb-3305-4ed9-b883-4623cbeed2f6"],"resources":["org:f5528d02-e6f8-40c6-9dd6-4531e96867ed:RBAC:permission","org:f5528d02-e6f8-40c6-9dd6-4531e96867ed:RBAC:role"],"actions":["create","view","assign","delete","grantPermission","delegatePermission"],"effect":"allow","conditions":null}'), 
    (2, '/store/ory/exact/policies', 'ac1fa1c8-bf00-46bf-b99f-2834e6a42fa3', '{"id":"ac1fa1c8-bf00-46bf-b99f-2834e6a42fa3","description":"Create System Accounts","subjects":["78ffaeeb-3305-4ed9-b883-4623cbeed2f6","817b2362-0aa0-4819-a323-2001f1984d3f"],"resources":["org:f5528d02-e6f8-40c6-9dd6-4531e96867ed:ACCOUNT:org","org:f5528d02-e6f8-40c6-9dd6-4531e96867ed:ACCOUNT:user","org:f5528d02-e6f8-40c6-9dd6-4531e96867ed:ACCOUNT:service"],"actions":["create","invite","grantPermission","delegatePermission"],"effect":"allow","conditions":null}'), 
    (3, '/store/ory/exact/roles', '78ffaeeb-3305-4ed9-b883-4623cbeed2f6', '{"id":"78ffaeeb-3305-4ed9-b883-4623cbeed2f6","description":"","members":["59SUgTDANlU0MxvHeOOAYmO9"]}');