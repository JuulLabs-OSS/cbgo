#import <Foundation/Foundation.h>
#import <CoreBluetooth/CoreBluetooth.h>
#import "bt.h"

@implementation BTDlg

- (id)
init
{
    self = [super init];
    if (self == nil) {
        return nil;
    }

    return self;
}

/**
 * Called when the central manager successfully connects to a peripheral.
 */
- (void)
      centralManager:(CBCentralManager *)cm
didConnectPeripheral:(CBPeripheral *)prph
{
    prph.delegate = self;
    BTCentralManagerDidConnectPeripheral(cm, prph);
}

/**
 * Called when a connection to a preipheral is terminated.
 */
- (void)
         centralManager:(CBCentralManager *)cm
didDisconnectPeripheral:(CBPeripheral *)prph
                  error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTCentralManagerDidDisconnectPeripheral(cm, prph, &err);
}

/**
 * Called when the central manager fails to connect to a peripheral.
 */
- (void)
            centralManager:(CBCentralManager *)cm
didFailToConnectPeripheral:(CBPeripheral *)prph
                     error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTCentralManagerDidFailToConnectPeripheral(cm, prph, &err);
}

// macOS 10.15+
#if 0
- (void)
         centralManager:(CBCentralManager *)cm 
connectionEventDidOccur:(CBConnectionEvent)event 
          forPeripheral:(CBPeripheral *)prph
{
    BTCentralManagerConnectionEventDidOccur(cm, event, prph);
}
#endif

/**
 * Called when the central manager discovers a peripheral while scanning for
 * devices.
 */
- (void)
       centralManager:(CBCentralManager *)cm
didDiscoverPeripheral:(CBPeripheral *)prph
    advertisementData:(NSDictionary *)advData
                 RSSI:(NSNumber *)RSSI
{
    struct adv_fields af = {0};

    af.name = dict_string(advData, CBAdvertisementDataLocalNameKey);
    af.mfg_data = dict_bytes(advData, CBAdvertisementDataManufacturerDataKey);
    af.pwr_lvl = dict_int(advData, CBAdvertisementDataTxPowerLevelKey, ADV_FIELDS_PWR_LVL_NONE);
    af.connectable = dict_int(advData, CBAdvertisementDataIsConnectable, ADV_FIELDS_CONNECTABLE_NONE);

    const NSArray *arr = [advData objectForKey:CBAdvertisementDataServiceUUIDsKey];
    const char *svc_uuids[[arr count]];
    for (int i = 0; i < [arr count]; i++) {
        const CBUUID *uuid = [arr objectAtIndex:i];
        svc_uuids[i] = [[uuid UUIDString] UTF8String];
    }
    af.svc_uuids = (struct string_arr) {
        .strings = svc_uuids,
        .count = [arr count],
    };

    const NSDictionary *dict = [advData objectForKey:CBAdvertisementDataServiceDataKey];
    const NSArray *keys = [dict allKeys];

    const char *svc_data_uuids[[keys count]];
    struct byte_arr svc_data_values[[keys count]];

    for (int i = 0; i < [keys count]; i++) {
        const CBUUID *uuid = [keys objectAtIndex:i];
        svc_data_uuids[i] = [[uuid UUIDString] UTF8String];

        const NSData *data = [dict objectForKey:uuid];
        svc_data_values[i].data = [data bytes];
        svc_data_values[i].length = [data length];
    }
    af.svc_data_uuids = (struct string_arr) {
        .strings = svc_data_uuids,
        .count = [keys count],
    };
    af.svc_data_values = svc_data_values;

    prph.delegate = self;
    [prph retain];

    BTCentralManagerDidDiscoverPeripheral(cm, prph, &af, [RSSI intValue]);
}

/**
 * Called whenever the central manager's state is updated.
 */
- (void)
centralManagerDidUpdateState:(CBCentralManager *)cm
{
    BTCentralManagerDidUpdateState(cm);
}

- (void)
  centralManager:(CBCentralManager *)cm 
willRestoreState:(NSDictionary<NSString *,id> *)dict
{
    struct restore_opts opts = {0};

    const NSArray *prphs = [dict objectForKey:CBCentralManagerRestoredStatePeripheralsKey];
    opts.prphs = nsarray_to_obj_arr(prphs);

    const NSArray *uuids = [dict objectForKey:CBCentralManagerRestoredStateScanServicesKey];
    opts.scan_svcs = cbuuids_to_strs(uuids);

    struct scan_opts scan_opts = {0};
    const NSDictionary *scan_dict = [dict objectForKey:CBCentralManagerRestoredStateScanOptionsKey];
    if (scan_dict != nil && [scan_dict count] > 0) {
        opts.scan_opts = &scan_opts;

        NSNumber *dups = [scan_dict objectForKey:CBCentralManagerScanOptionAllowDuplicatesKey];
        if (dups != nil && [dups boolValue]) {
            opts.scan_opts->allow_dups = true;
        }

        NSArray *sol_uuids = [scan_dict objectForKey:CBCentralManagerScanOptionSolicitedServiceUUIDsKey];
        opts.scan_opts->sol_svc_uuids = cbuuids_to_strs(sol_uuids);
    }

    BTCentralManagerWillRestoreState(cm, &opts);

    free(scan_opts.sol_svc_uuids.strings);
    free(opts.scan_svcs.strings);
    free(opts.prphs.objs);
}

/**
 * Called when the central manager successfully discovers services on a
 * peripheral.
 */
- (void) peripheral:(CBPeripheral *)prph
didDiscoverServices:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidDiscoverServices(prph, &err);
}

- (void)                   peripheral:(CBPeripheral *)prph
didDiscoverIncludedServicesForService:(CBService *)svc
                                error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidDiscoverIncludedServices(prph, svc, &err);
}

- (void)
                          peripheral:(CBPeripheral *)prph 
didDiscoverCharacteristicsForService:(CBService *)svc 
                               error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidDiscoverCharacteristics(prph, svc, &err);
}

- (void)
                             peripheral:(CBPeripheral *)prph 
didDiscoverDescriptorsForCharacteristic:(CBCharacteristic *)chr 
             error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidDiscoverDescriptors(prph, chr, &err);
}

- (void)
                     peripheral:(CBPeripheral *)prph 
didUpdateValueForCharacteristic:(CBCharacteristic *)chr 
                          error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidUpdateValueForCharacteristic(prph, chr, &err);
}

- (void)
                 peripheral:(CBPeripheral *)prph 
didUpdateValueForDescriptor:(CBDescriptor *)dsc 
                      error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidUpdateValueForDescriptor(prph, dsc, &err);
}

- (void)
                    peripheral:(CBPeripheral *)prph 
didWriteValueForCharacteristic:(CBCharacteristic *)chr 
                         error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidWriteValueForCharacteristic(prph, chr, &err);
}

- (void)
                peripheral:(CBPeripheral *)prph 
didWriteValueForDescriptor:(CBDescriptor *)dsc 
                     error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidWriteValueForDescriptor(prph, dsc, &err);
}

- (void)
peripheralIsReadyToSendWriteWithoutResponse:(CBPeripheral *)prph
{
    BTPeripheralIsReadyToSendWriteWithoutResponse(prph);
}

- (void)
                                 peripheral:(CBPeripheral *)prph 
didUpdateNotificationStateForCharacteristic:(CBCharacteristic *)chr 
                                      error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidUpdateNotificationState(prph, chr, &err);
}

- (void)
 peripheral:(CBPeripheral *)prph 
didReadRSSI:(NSNumber *)RSSI 
      error:(NSError *)nserr
{
    struct bt_error err = nserror_to_bt_error(nserr);
    BTPeripheralDidReadRSSI(prph, [RSSI intValue], &err);
}

- (void)
peripheralDidUpdateName:(CBPeripheral *)prph
{
    BTPeripheralDidUpdateName(prph);
}

- (void)
       peripheral:(CBPeripheral *)prph 
didModifyServices:(NSArray<CBService *> *)invSvcs
{
    struct obj_arr oa = nsarray_to_obj_arr(invSvcs);
    BTPeripheralDidModifyServices(prph, &oa);
    free(oa.objs);
}

@end
