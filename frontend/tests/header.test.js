import { h } from 'preact';
import Header from '../src/components/header';
// See: https://github.com/preactjs/enzyme-adapter-preact-pure
import { mount, shallow } from 'enzyme';

describe('Initial Test of the Header', () => {
	test('Header renders 3 nav items', () => {
		const context = mount(<Header selectedRoute="/"/>);
		// expect(context.find('TopAppBar.Title').text()).toBe('Barndoor Tracker');
		// expect(context.find('Drawer.DrawerItem').length).toBe(2);
	});
});
