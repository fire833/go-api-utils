/*
*	Copyright (C) 2024 Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package mgr

// AddItem appends a new SubsystemStatus object to the existing list of items within the existing SubsystemStatusList.
func (a *SubsystemStatusList) AddItem(item *SubsystemStatus) {
	a.Items = append(a.Items, item)
}

// AddItem appends a new BuildInfo object to the existing list of items within the existing BuildInfoList.
func (a *BuildInfoList) AddItem(item *BuildInfo) {
	a.Items = append(a.Items, item)
}
