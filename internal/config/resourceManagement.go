/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package config

var resources chan struct{}

// GetResourceTracker returns a buffered channel to act as a count of the active connections to subsystems and block when too many connections are open
func GetResourceTracker() chan struct{} {
	if V.ActiveFileProcessing == 0 {
		resource := make(chan struct{}, 1)
		resource <- struct{}{}
		return resource
	}
	return resources
}
