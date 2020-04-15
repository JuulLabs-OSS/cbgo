#ifndef H_BT_
#define H_BT_

#import <Foundation/Foundation.h>
#import <CoreBluetooth/CoreBluetooth.h>

#define ADV_FIELDS_PWR_LVL_NONE         (-128)
#define ADV_FIELDS_CONNECTABLE_NONE     (-1)

struct byte_arr {
    const uint8_t *data;
    int length;
};

struct string_arr {
    const char **strings;
    int count;
};

struct obj_arr {
    void **objs;
    int count;
};

struct bt_error {
    const char *msg;
    int code;
};

struct adv_fields {
	const char *name;
	struct byte_arr mfg_data;
    struct string_arr svc_uuids;
    struct string_arr overflow_svc_uuids;
    int8_t pwr_lvl; // -128 = not present
	int connectable; // -1 = not present
    struct string_arr svc_data_uuids;
    struct byte_arr *svc_data_values;
};

struct scan_opts {
    bool allow_dups;
    struct string_arr sol_svc_uuids;
};

struct connect_opts {
	bool notify_on_connection;
	bool notify_on_disconnection;
	bool notify_on_notification;
	bool enable_transport_bridging;
	bool requires_ancs;
	int start_delay;
};

struct restore_opts {
    struct obj_arr prphs;
    struct string_arr scan_svcs;
    struct scan_opts *scan_opts;
};

@interface BTDlg : NSObject <CBCentralManagerDelegate, CBPeripheralDelegate>
{
}
@end

// bt.m
bool bt_start();
void bt_stop();
void bt_init();

// util.m
struct byte_arr nsdata_to_byte_arr(const NSData *nsdata);
NSData *byte_arr_to_nsdata(const struct byte_arr *ba);
struct obj_arr nsarray_to_obj_arr(const NSArray *arr);
NSString *str_to_nsstring(const char *s);
struct bt_error nserror_to_bt_error(const NSError *err);
struct string_arr cbuuids_to_strs(const NSArray *cbuuids);
int dict_int(NSDictionary *dict, NSString *key, int dflt);
const char *dict_string(NSDictionary *dict, NSString *key);
const void *dict_data(NSDictionary *dict, NSString *key, int *out_len);
const struct byte_arr dict_bytes(NSDictionary *dict, NSString *key);
void dict_set_bool(NSMutableDictionary *dict, NSString *key, bool val);
void dict_set_int(NSMutableDictionary *dict, NSString *key, int val);
NSUUID *str_to_nsuuid(const char *s);
CBUUID *str_to_cbuuid(const char *s);
NSArray *strs_to_nsuuids(const struct string_arr *sa);
NSArray *strs_to_cbuuids(const struct string_arr *sa);
NSArray *strs_to_nsstrings(const struct string_arr *sa);

// cb.m
CBCentralManager *cb_alloc_cmgr(bool pwr_alert, const char *restore_id);
void cb_cmgr_set_delegate(void *cmgr, bool set);
int cb_cmgr_state(void *cm);
void cb_cmgr_scan(void *cmgr, const struct string_arr *svc_uuids,
                  const struct scan_opts *opts);
void cb_cmgr_stop_scan(void *cm);
bool cb_cmgr_is_scanning(void *cm);
void cb_cmgr_connect(void *cmgr, void *prph, const struct connect_opts *opts);
void cb_cmgr_cancel_connect(void *cmgr, void *prph);
struct obj_arr cb_cmgr_retrieve_prphs_with_svcs(void *cmgr, const struct string_arr *svc_uuids);
struct obj_arr cb_cmgr_retrieve_prphs(void *cmgr, const struct string_arr *uuids);

void cb_prph_set_delegate(void *prph, bool set);
const char *cb_prph_identifier(void *prph);
const char *cb_prph_name(void *prph);
struct obj_arr cb_prph_services(void *prph);
void cb_prph_discover_svcs(void *prph, const struct string_arr *svc_uuid_strs);
void cb_prph_discover_included_svcs(void *prph, const struct string_arr *svc_uuid_strs, void *svc);
void cb_prph_discover_chrs(void *prph, void *svc, const struct string_arr *chr_uuid_strs);
void cb_prph_discover_dscs(void *prph, void *chr);
void cb_prph_read_chr(void *prph, void *chr);
void cb_prph_read_dsc(void *prph, void *dsc);
void cb_prph_write_chr(void *prph, void *chr, struct byte_arr *value, int type);
void cb_prph_write_dsc(void *prph, void *dsc, struct byte_arr *value);
int cb_prph_max_write_len(void *prph, int type);
void cb_prph_set_notify(void *prph, bool enabled, void *chr);
int cb_prph_state(void *prph);
bool cb_prph_can_send_write_without_rsp(void *prph);
void cb_prph_read_rssi(void *prph);
bool cb_prph_ancs_authorized(void *prph);

const char *cb_svc_uuid(void *svc);
void *cb_svc_peripheral(void *svc);
bool cb_svc_is_primary(void *svc);
struct obj_arr cb_svc_characteristics(void *svc);
struct obj_arr cb_svc_included_svcs(void *svc);

const char *cb_chr_uuid(void *chr);
void *cb_chr_service(void *chr);
struct obj_arr cb_chr_descriptors(void *chr);
struct byte_arr cb_chr_value(void *chr);
int cb_chr_properties(void *chr);
bool cb_chr_is_notifying(void *chr);

const char *cb_dsc_uuid(void *dsc);
void *cb_dsc_characteristic(void *dsc);
struct byte_arr cb_dsc_value(void *dsc);

// cbhandlers.go
void BTCentralManagerDidConnectPeripheral(void *cmgr, void *prph);
void BTCentralManagerDidFailToConnectPeripheral(void *cmgr, void *prph, struct bt_error *err);
void BTCentralManagerDidDisconnectPeripheral(void *cmgr, void *prph, struct bt_error *err);
void BTCentralManagerConnectionEventDidOccur(void *cmgr, int event, void *prph);
void BTCentralManagerDidDiscoverPeripheral(void *cmgr, void *prph, struct adv_fields *advData, int rssi);
void BTCentralManagerDidUpdateState(void *cmgr);
void BTCentralManagerWillRestoreState(void *cmgr, struct restore_opts *opts);
void BTPeripheralDidDiscoverServices(void *prph, struct bt_error *err);
void BTPeripheralDidDiscoverIncludedServices(void *prph, void *svc, struct bt_error *err);
void BTPeripheralDidDiscoverCharacteristics(void *prph, void *svc, struct bt_error *err);
void BTPeripheralDidDiscoverDescriptors(void *prph, void *chr, struct bt_error *err);
void BTPeripheralDidUpdateValueForCharacteristic(void *prph, void *chr, struct bt_error *err);
void BTPeripheralDidUpdateValueForDescriptor(void *prph, void *dsc, struct bt_error *err);
void BTPeripheralDidWriteValueForCharacteristic(void *prph, void *chr, struct bt_error *err);
void BTPeripheralDidWriteValueForDescriptor(void *prph, void *dsc, struct bt_error *err);
void BTPeripheralIsReadyToSendWriteWithoutResponse(void *prph);
void BTPeripheralDidUpdateNotificationState(void *prph, void *chr, struct bt_error *err);
void BTPeripheralDidReadRSSI(void *prph, int rssi, struct bt_error *err);
void BTPeripheralDidUpdateName(void *prph);
void BTPeripheralDidModifyServices(void *prph, struct obj_arr *inv_svcs);

extern dispatch_queue_t bt_queue;
extern BTDlg *bt_dlg;

#endif
