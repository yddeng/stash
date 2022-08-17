
DROP TABLE IF EXISTS "rank_list";
CREATE TABLE "rank_list" (
      "__key__" varchar(255) NOT NULL,
      "__version__" int8 NOT NULL,
      "rank_info" bytea NOT NULL DEFAULT '',
      "rank_status" int8 NOT NULL DEFAULT 0,
      "settled_idx" int8 NOT NULL DEFAULT 0,
      "logic_addr" varchar(255) NOT NULL DEFAULT '',
      "last_id" int8 NOT NULL DEFAULT 0,
      PRIMARY KEY ("__key__")
);

/*    ****************************************************       */
/*    ********************  *****************************       */
/*    ****************************************************       */
DROP TABLE IF EXISTS "game_user_0";
DROP TABLE IF EXISTS "user_module_data_0";
DROP TABLE IF EXISTS "user_game_login_0";
DROP TABLE IF EXISTS "role_name_0";
DROP TABLE IF EXISTS "character_0";
DROP TABLE IF EXISTS "user_dir_last_login_0";
DROP TABLE IF EXISTS "global_data_0";
DROP TABLE IF EXISTS "function_switch_0";
DROP TABLE IF EXISTS "whitelist_0";
DROP TABLE IF EXISTS "user_assets_0";
DROP TABLE IF EXISTS "backpack_0";
DROP TABLE IF EXISTS "equip_0";
DROP TABLE IF EXISTS "weapon_0";
DROP TABLE IF EXISTS "msgbox_0";
DROP TABLE IF EXISTS "main_dungeons_0";
DROP TABLE IF EXISTS "quest_0";
DROP TABLE IF EXISTS "drawcard_0";
DROP TABLE IF EXISTS "scarsingrain_0";
DROP TABLE IF EXISTS "materialdungeon_0";
DROP TABLE IF EXISTS "rewardquest_0";
DROP TABLE IF EXISTS "temporary_0";
DROP TABLE IF EXISTS "nodelock_0";
DROP TABLE IF EXISTS "firewall_0";
DROP TABLE IF EXISTS "whitelist_0";
DROP TABLE IF EXISTS "bannedlist_0";
DROP TABLE IF EXISTS "shop_0";
DROP TABLE IF EXISTS "worldquest_0";
DELETE FROM "rank_list";
DROP TABLE IF EXISTS "id_counter_0";
DROP TABLE IF EXISTS "mail_0";
DROP TABLE IF EXISTS "offlinemsg_0";
DROP TABLE IF EXISTS "sign_0";
DROP TABLE IF EXISTS "web_data_0";
